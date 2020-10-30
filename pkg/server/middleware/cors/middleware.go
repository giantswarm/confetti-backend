package cors

import (
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

func (m *Middleware) Middleware() atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", m.flags.AllowedOrigin)

		return ctx.Next()
	}
}
