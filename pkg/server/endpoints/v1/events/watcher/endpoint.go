package watcher

import (
	"github.com/atreugo/websocket"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
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
	Hub               *websocketutil.Hub
}

type Endpoint struct {
	flags             *flags.Flags
	service           *Service
	middleware        *middleware.Middleware
	websocketUpgrader *websocket.Upgrader
	hub               *websocketutil.Hub
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
	if c.Hub == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Hub must not be empty", c)
	}

	go c.Hub.Run()

	endpoint := &Endpoint{
		flags:             c.Flags,
		service:           c.Service,
		middleware:        c.Middleware,
		websocketUpgrader: c.WebsocketUpgrader,
		hub:               c.Hub,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return e.websocketUpgrader.Upgrade(func(ws *websocket.Conn) error {
		c := websocketutil.ClientConfig{
			Hub:        e.hub,
			Connection: ws,
		}

		client, err := websocketutil.NewClient(c)
		if err != nil {
			return microerror.Mask(err)
		}

		go client.WritePump()
		client.ReadPump()

		return nil
	})
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{
		Before: []atreugo.Middleware{
			e.middleware.Authentication.Middleware,
			ValidateIDMiddleware(e.service),
		},
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
