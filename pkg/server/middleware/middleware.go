package middleware

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware/cors"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware/users"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type Config struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Middleware struct {
	Users *users.Middleware
	Cors  *cors.Middleware

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

	var usersMiddleware *users.Middleware
	{
		c := users.Config{
			Flags:  c.Flags,
			Models: c.Models,
		}

		usersMiddleware, err = users.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var corsMiddleware *cors.Middleware
	{
		c := cors.Config{
			Flags:  c.Flags,
			Models: c.Models,
		}

		corsMiddleware, err = cors.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	m := &Middleware{
		Users: usersMiddleware,
		Cors:  corsMiddleware,

		flags:  c.Flags,
		models: c.Models,
	}

	return m, nil
}
