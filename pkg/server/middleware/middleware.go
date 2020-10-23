package middleware

import (
	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware/authentication"
	"github.com/giantswarm/microerror"
)

type Config struct {
	Flags *flags.Flags
}

type Middleware struct {
	Authentication *authentication.Middleware

	flags *flags.Flags
}

func New(c Config) (*Middleware, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}

	var err error

	var authenticationMiddleware *authentication.Middleware
	{
		c := authentication.Config{
			Flags: c.Flags,
		}

		authenticationMiddleware, err = authentication.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	m := &Middleware{
		Authentication: authenticationMiddleware,

		flags: c.Flags,
	}

	return m, nil
}