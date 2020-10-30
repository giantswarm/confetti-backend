package events

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware/events/validateid"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type Config struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Middleware struct {
	ValidateID *validateid.Middleware

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

	var validateidMiddleware *validateid.Middleware
	{
		c := validateid.Config{
			Flags:  c.Flags,
			Models: c.Models,
		}

		validateidMiddleware, err = validateid.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	m := &Middleware{
		ValidateID: validateidMiddleware,

		flags:  c.Flags,
		models: c.Models,
	}

	return m, nil
}
