package users

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/users/login"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type EndpointConfig struct {
	Flags      *flags.Flags
	Middleware *middleware.Middleware
	Models     *models.Model
}

type Endpoint struct {
	Login *login.Endpoint

	flags      *flags.Flags
	middleware *middleware.Middleware
	models     *models.Model
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Middleware must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}

	loginEndpoint, err := createLoginEndpoint(c.Flags, c.Middleware, c.Models)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Login: loginEndpoint,

		flags:      c.Flags,
		middleware: c.Middleware,
		models:     c.Models,
	}

	return endpoint, nil
}

func createLoginEndpoint(flags *flags.Flags, middleware *middleware.Middleware, models *models.Model) (*login.Endpoint, error) {
	var err error

	var service *login.Service
	{
		c := login.ServiceConfig{
			Flags:  flags,
			Models: models,
		}
		service, err = login.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoint *login.Endpoint
	{
		c := login.EndpointConfig{
			Flags:      flags,
			Service:    service,
			Middleware: middleware,
			Models:     models,
		}
		endpoint, err = login.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
