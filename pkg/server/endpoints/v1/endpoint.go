package v1

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events"
	eventsModel "github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/users"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
)

const (
	method = "GET"
	path   = "/"
)

type EndpointConfig struct {
	Flags      *flags.Flags
	Middleware *middleware.Middleware
}

type Endpoint struct {
	Users *users.Endpoint
	Events *events.Endpoint

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

	usersEndpoint, err := createUsersEndpoint(c.Flags, c.Middleware)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	eventsEndpoint, err := createEventsEndpoint(c.Flags, c.Middleware)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Users: usersEndpoint,
		Events: eventsEndpoint,

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

func createUsersEndpoint(flags *flags.Flags, middleware *middleware.Middleware) (*users.Endpoint, error) {
	var err error

	var endpoint *users.Endpoint
	{
		c := users.EndpointConfig{
			Flags:      flags,
			Middleware: middleware,
		}
		endpoint, err = users.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func createEventsEndpoint(flags *flags.Flags, middleware *middleware.Middleware) (*events.Endpoint, error) {
	var err error

	var repository *eventsModel.Repository
	{
		c := eventsModel.RepositoryConfig{
			Flags: flags,
		}
		repository, err = eventsModel.NewRepository(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var service *events.Service
	{
		c := events.ServiceConfig{
			Flags: flags,
			Repository: repository,
		}
		service, err = events.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoint *events.Endpoint
	{
		c := events.EndpointConfig{
			Flags:      flags,
			Service:    service,
			Middleware: middleware,
		}
		endpoint, err = events.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
