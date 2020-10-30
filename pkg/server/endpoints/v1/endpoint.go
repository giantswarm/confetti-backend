package v1

import (
	"net/http"

	"github.com/atreugo/websocket"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/users"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

const (
	method = "GET"
	path   = "/"
)

type EndpointConfig struct {
	Flags             *flags.Flags
	Middleware        *middleware.Middleware
	Models            *models.Model
	WebsocketUpgrader *websocket.Upgrader
}

type Endpoint struct {
	Users  *users.Endpoint
	Events *events.Endpoint

	flags             *flags.Flags
	middleware        *middleware.Middleware
	models            *models.Model
	websocketUpgrader *websocket.Upgrader
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Middleware must not be empty", c)
	}
	if c.WebsocketUpgrader == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.WebsocketUpgrader must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}

	usersEndpoint, err := createUsersEndpoint(c.Flags, c.Middleware, c.Models)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	eventsEndpoint, err := createEventsEndpoint(c.Flags, c.Middleware, c.Models, c.WebsocketUpgrader)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Users:  usersEndpoint,
		Events: eventsEndpoint,

		flags:             c.Flags,
		middleware:        c.Middleware,
		models:            c.Models,
		websocketUpgrader: c.WebsocketUpgrader,
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

func createUsersEndpoint(flags *flags.Flags, middleware *middleware.Middleware, models *models.Model) (*users.Endpoint, error) {
	var err error

	var endpoint *users.Endpoint
	{
		c := users.EndpointConfig{
			Flags:      flags,
			Middleware: middleware,
			Models:     models,
		}
		endpoint, err = users.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func createEventsEndpoint(flags *flags.Flags, middleware *middleware.Middleware, models *models.Model, websocketUpgrader *websocket.Upgrader) (*events.Endpoint, error) {
	var err error

	var service *events.Service
	{
		c := events.ServiceConfig{
			Flags:  flags,
			Models: models,
		}
		service, err = events.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoint *events.Endpoint
	{
		c := events.EndpointConfig{
			Flags:             flags,
			Service:           service,
			Middleware:        middleware,
			Models:            models,
			WebsocketUpgrader: websocketUpgrader,
		}
		endpoint, err = events.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
