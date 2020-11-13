package validateid

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/context/event"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	eventsModel "github.com/giantswarm/confetti-backend/pkg/server/models/events"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type Config struct {
	Flags  *flags.Flags
	Models *models.Model
	Logger micrologger.Logger
}

type Middleware struct {
	flags  *flags.Flags
	models *models.Model
	logger micrologger.Logger
}

func New(c Config) (*Middleware, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", c)
	}

	m := &Middleware{
		flags:  c.Flags,
		models: c.Models,
		logger: c.Logger,
	}

	return m, nil
}

// Middleware validates if a event with the event ID
// URL param exists.
func (m *Middleware) Middleware() atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		id, ok := ctx.UserValue("id").(string)
		if !ok {
			ctx.SetStatusCode(http.StatusBadRequest)

			return microerror.Mask(invalidParamsError)
		}

		e, err := m.findEventByID(id)
		if eventsModel.IsNotFoundError(err) {
			ctx.SetStatusCode(http.StatusNotFound)

			return microerror.Mask(notFoundError)
		} else if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)

			return microerror.Mask(internalServerError)
		}
		event.SaveContext(ctx, e)

		return ctx.Next()
	}
}

func (m *Middleware) findEventByID(id string) (eventsModelTypes.Event, error) {
	e, err := m.models.Events.FindOneByID(id)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return e, nil
}
