package watcher

import (
	"fmt"

	"github.com/atreugo/websocket"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
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
}

type Endpoint struct {
	flags             *flags.Flags
	service           *Service
	middleware        *middleware.Middleware
	websocketUpgrader *websocket.Upgrader
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

	endpoint := &Endpoint{
		flags:             c.Flags,
		service:           c.Service,
		middleware:        c.Middleware,
		websocketUpgrader: c.WebsocketUpgrader,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return e.websocketUpgrader.Upgrade(func(ws *websocket.Conn) error {
		for {
			msg := make(map[string]string)
			err := ws.ReadJSON(&msg)
			if err != nil {
				return microerror.Mask(err)
			}

			fmt.Printf("recv: %s\n", msg)

			err = ws.WriteJSON(msg)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	})
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{
		// Before: []atreugo.Middleware{
		// 	e.middleware.Authentication.Middleware,
		// },
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
