package events

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
)

const (
	method = "GET"
	path   = "/events/"
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
		events, err := e.service.GetEvents()
		if err != nil {
			rErr := ctx.ErrorResponse(microerror.Mask(err), http.StatusInternalServerError)
			if rErr != nil {
				return microerror.Mask(err)
			}
		}

		res := Response{}
		{
			res.Events = make([]ResponseEvent, 0, len(events))
			for _, e := range events {
				res.Events = append(res.Events, ResponseEvent{
					Active: e.Active(),
					ID: e.ID(),
					Name: e.Name(),
					EventType : string(e.Type()),
				})
			}
		}


		return ctx.JSONResponse(res, http.StatusOK)
	}
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{
		Before: []atreugo.Middleware{
			e.middleware.Authentication.Middleware,
		},
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
