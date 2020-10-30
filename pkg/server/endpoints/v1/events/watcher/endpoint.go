package watcher

import (
	"github.com/atreugo/websocket"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

const (
	method = "GET"
	path   = "/events/{id}/watch/"
)

type EndpointConfig struct {
	Flags             *flags.Flags
	Service           *Service
	Middleware        *middleware.Middleware
	WebsocketUpgrader *websocket.Upgrader
	Models            *models.Model
}

type Endpoint struct {
	flags             *flags.Flags
	service           *Service
	middleware        *middleware.Middleware
	websocketUpgrader *websocket.Upgrader
	models            *models.Model
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Service == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Service must not be empty", c)
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

	endpoint := &Endpoint{
		flags:             c.Flags,
		service:           c.Service,
		middleware:        c.Middleware,
		models:            c.Models,
		websocketUpgrader: c.WebsocketUpgrader,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return e.websocketUpgrader.Upgrade(func(ws *websocket.Conn) error {
		err := e.service.NewClient(ws)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	})
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{
		Before: []atreugo.Middleware{
			e.middleware.Users.Authentication.Middleware(),
			e.middleware.Events.ValidateID.Middleware(),
		},
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
