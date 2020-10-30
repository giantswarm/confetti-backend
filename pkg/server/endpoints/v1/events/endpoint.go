package events

import (
	"net/http"

	"github.com/atreugo/websocket"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/searcher"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
)

const (
	method = "GET"
	path   = "/events/"
)

type EndpointConfig struct {
	Flags             *flags.Flags
	Service           *Service
	Middleware        *middleware.Middleware
	WebsocketUpgrader *websocket.Upgrader
}

type Endpoint struct {
	Searcher *searcher.Endpoint
	Watcher  *watcher.Endpoint

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

	searcherEndpoint, err := createSearcherEndpoint(c.Flags, c.Middleware, c.Service.repository)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	watcherEndpoint, err := createWatcherEndpoint(c.Flags, c.Middleware, c.Service.repository, c.WebsocketUpgrader)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Searcher: searcherEndpoint,
		Watcher:  watcherEndpoint,

		flags:             c.Flags,
		service:           c.Service,
		middleware:        c.Middleware,
		websocketUpgrader: c.WebsocketUpgrader,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		events, err := e.service.GetEvents()
		if err != nil {
			return ctx.ErrorResponse(microerror.Mask(err), http.StatusInternalServerError)
		}

		res := Response{}
		{
			res.Events = make([]ResponseEvent, 0, len(events))
			for _, e := range events {
				res.Events = append(res.Events, ResponseEvent{
					Active:    e.Active(),
					ID:        e.ID(),
					Name:      e.Name(),
					EventType: string(e.Type()),
				})
			}
		}

		return ctx.JSONResponse(res, http.StatusOK)
	}
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{
		Before: []atreugo.Middleware{
			e.middleware.Authentication.Middleware,
		},
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}

func createSearcherEndpoint(flags *flags.Flags, middleware *middleware.Middleware, repository *model.Repository) (*searcher.Endpoint, error) {
	var err error

	var service *searcher.Service
	{
		c := searcher.ServiceConfig{
			Flags:      flags,
			Repository: repository,
		}
		service, err = searcher.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoint *searcher.Endpoint
	{
		c := searcher.EndpointConfig{
			Flags:      flags,
			Service:    service,
			Middleware: middleware,
		}
		endpoint, err = searcher.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func createWatcherEndpoint(flags *flags.Flags, middleware *middleware.Middleware, repository *model.Repository, websocketUpgrader *websocket.Upgrader) (*watcher.Endpoint, error) {
	var err error

	var hub *websocketutil.Hub
	{
		hub, err = websocketutil.NewHub()
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var service *watcher.Service
	{
		c := watcher.ServiceConfig{
			Flags:      flags,
			Repository: repository,
		}
		service, err = watcher.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoint *watcher.Endpoint
	{
		c := watcher.EndpointConfig{
			Flags:             flags,
			Service:           service,
			Middleware:        middleware,
			WebsocketUpgrader: websocketUpgrader,
			Hub:               hub,
		}
		endpoint, err = watcher.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
