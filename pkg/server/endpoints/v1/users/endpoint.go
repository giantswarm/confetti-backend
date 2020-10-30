package users

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/users/login"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
)

const (
	method = "GET"
	path   = "/users/"
)

type EndpointConfig struct {
	Flags      *flags.Flags
	Middleware *middleware.Middleware
}

type Endpoint struct {
	Login *login.Endpoint

	flags      *flags.Flags
	middleware *middleware.Middleware
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Middleware must not be empty", c)
	}

	loginEndpoint, err := createLoginEndpoint(c.Flags, c.Middleware)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Login: loginEndpoint,

		flags:      c.Flags,
		middleware: c.Middleware,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		err := ctx.HTTPResponse("", http.StatusNotFound)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}

func createLoginEndpoint(flags *flags.Flags, middleware *middleware.Middleware) (*login.Endpoint, error) {
	var err error

	var service *login.Service
	{
		c := login.ServiceConfig{
			Flags: flags,
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
		}
		endpoint, err = login.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
