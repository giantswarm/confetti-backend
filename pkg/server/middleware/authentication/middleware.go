package authentication

import (
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flags"
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
	err := ctx.Next()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
