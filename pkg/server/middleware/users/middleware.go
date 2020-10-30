package users

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware/users/authentication"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type Config struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Middleware struct {
	Authentication *authentication.Middleware

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

	var err error

	var authenticationMiddleware *authentication.Middleware
	{
		c := authentication.Config{
			Flags:  c.Flags,
			Models: c.Models,
		}

		authenticationMiddleware, err = authentication.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	m := &Middleware{
		Authentication: authenticationMiddleware,

		flags:  c.Flags,
		models: c.Models,
	}

	return m, nil
}
