package cors

import (
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/internal/flags"
)

type Config struct {
	Flags *flags.Flags
}

type Middleware struct {
	flags *flags.Flags
}

func New(c Config) (*Middleware, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}

	m := &Middleware{
		flags: c.Flags,
	}

	return m, nil
}

func (m *Middleware) Middleware(ctx *atreugo.RequestCtx) error {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", m.flags.AllowedOrigin)

	return ctx.Next()
}
