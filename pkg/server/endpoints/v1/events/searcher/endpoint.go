package searcher

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/searcher/response"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	eventsModel "github.com/giantswarm/confetti-backend/pkg/server/models/events"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

const (
	method = "GET"
	path   = "/events/{id}"
)

type EndpointConfig struct {
	Flags      *flags.Flags
	Service    *Service
	Middleware *middleware.Middleware
}

type Endpoint struct {
	flags      *flags.Flags
	service    *Service
	middleware *middleware.Middleware
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Service == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Service must not be empty", c)
	}
	if c.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Middleware must not be empty", c)
	}

	endpoint := &Endpoint{
		flags:      c.Flags,
		service:    c.Service,
		middleware: c.Middleware,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		id, ok := ctx.UserValue("id").(string)
		if !ok {
			return ctx.ErrorResponse(microerror.Mask(invalidParamsError), http.StatusBadRequest)
		}

		event, err := e.service.GetEventByID(id)
		if eventsModel.IsNotFoundError(err) {
			return ctx.ErrorResponse(microerror.Mask(err), http.StatusNotFound)
		} else if err != nil {
			return ctx.ErrorResponse(microerror.Mask(err), http.StatusInternalServerError)
		}

		res := response.Response{
			Active:    event.Active(),
			ID:        event.ID(),
			Name:      event.Name(),
			EventType: string(event.Type()),
		}
		{
			// Add event type-specific details.
			switch e := event.(type) {
			case *eventsModelTypes.OnsiteEvent:
				res.Details.Rooms = make([]response.ResponseOnsiteRoom, 0, len(e.Rooms))
				for _, room := range e.Rooms {
					res.Details.Rooms = append(res.Details.Rooms, response.ResponseOnsiteRoom{
						ID:            room.ID,
						Name:          room.Name,
						Description:   room.Description,
						ConferenceURL: room.ConferenceURL,
					})
				}
			}
		}

		return ctx.JSONResponse(res, http.StatusOK)
	}
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{
		Before: []atreugo.Middleware{
			e.middleware.Authentication.Middleware(),
		},
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
