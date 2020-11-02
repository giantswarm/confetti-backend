package validateid

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type Config struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Middleware struct {
	flags  *flags.Flags
	models *models.Model
}

func New(c Config) (*Middleware, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}

	m := &Middleware{
		flags:  c.Flags,
		models: c.Models,
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

		if !m.eventExists(id) {
			ctx.SetStatusCode(http.StatusNotFound)

			return microerror.Mask(notFoundError)
		}

		return ctx.Next()
	}
}

func (m *Middleware) eventExists(id string) bool {
	_, err := m.models.Events.FindOneByID(id)

	return err == nil
}
