package events

import (
	"github.com/atreugo/websocket"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/lister"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/searcher"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
)

type EndpointConfig struct {
	Flags             *flags.Flags
	Middleware        *middleware.Middleware
	Models            *models.Model
	WebsocketUpgrader *websocket.Upgrader
}

type Endpoint struct {
	Searcher *searcher.Endpoint
	Watcher  *watcher.Endpoint
	Lister   *lister.Endpoint

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
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}
	if c.WebsocketUpgrader == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.WebsocketUpgrader must not be empty", c)
	}

	searcherEndpoint, err := createSearcherEndpoint(c.Flags, c.Middleware, c.Models)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	watcherEndpoint, err := createWatcherEndpoint(c.Flags, c.Middleware, c.Models, c.WebsocketUpgrader)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	listerEndpoint, err := createListerEndpoint(c.Flags, c.Middleware, c.Models)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Searcher: searcherEndpoint,
		Watcher:  watcherEndpoint,
		Lister:   listerEndpoint,

		flags:             c.Flags,
		middleware:        c.Middleware,
		models:            c.Models,
		websocketUpgrader: c.WebsocketUpgrader,
	}

	return endpoint, nil
}

func createSearcherEndpoint(flags *flags.Flags, middleware *middleware.Middleware, models *models.Model) (*searcher.Endpoint, error) {
	var err error

	var service *searcher.Service
	{
		c := searcher.ServiceConfig{
			Flags:  flags,
			Models: models,
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
			Models:     models,
		}
		endpoint, err = searcher.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func createWatcherEndpoint(flags *flags.Flags, middleware *middleware.Middleware, models *models.Model, websocketUpgrader *websocket.Upgrader) (*watcher.Endpoint, error) {
	var err error

	var hub websocketutil.Hub
	{
		hub, err = websocketutil.NewSocketHub()
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var service *watcher.Service
	{
		c := watcher.ServiceConfig{
			Flags:  flags,
			Models: models,
			Hub:    hub,
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
			Models:            models,
			WebsocketUpgrader: websocketUpgrader,
		}
		endpoint, err = watcher.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func createListerEndpoint(flags *flags.Flags, middleware *middleware.Middleware, models *models.Model) (*lister.Endpoint, error) {
	var err error

	var service *lister.Service
	{
		c := lister.ServiceConfig{
			Flags:  flags,
			Models: models,
		}
		service, err = lister.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoint *lister.Endpoint
	{
		c := lister.EndpointConfig{
			Flags:      flags,
			Service:    service,
			Middleware: middleware,
			Models:     models,
		}
		endpoint, err = lister.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
