package server

import (
	"github.com/atreugo/websocket"
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/root"
	v1 "github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type Config struct {
	Atreugo           *atreugo.Atreugo
	Flags             *flags.Flags
	WebsocketUpgrader *websocket.Upgrader
}

type Server struct {
	atreugo           *atreugo.Atreugo
	flags             *flags.Flags
	websocketUpgrader *websocket.Upgrader
}

func New(c Config) (*Server, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Atreugo == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Atreugo must not be empty", c)
	}
	if c.WebsocketUpgrader == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.WebsocketUpgrader must not be empty", c)
	}

	var err error

	s := &Server{
		atreugo:           c.Atreugo,
		flags:             c.Flags,
		websocketUpgrader: c.WebsocketUpgrader,
	}

	var allModels *models.Model
	{
		c := models.Config{
			Flags: s.flags,
		}
		allModels, err = models.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var allMiddlewares *middleware.Middleware
	{
		c := middleware.Config{
			Flags:  s.flags,
			Models: allModels,
		}
		allMiddlewares, err = middleware.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s.atreugo.UseBefore(allMiddlewares.Cors.Middleware())

	var rootEndpoint *root.Endpoint
	{
		rootEndpoint, err = newRootEndpoint(s.flags, allMiddlewares, allModels)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		s.atreugo.Path(
			rootEndpoint.Method(),
			rootEndpoint.Path(),
			rootEndpoint.Endpoint(),
		).Middlewares(rootEndpoint.Middlewares())
	}

	var v1Endpoint *v1.Endpoint
	{
		group := s.atreugo.NewGroupPath("/v1")

		v1Endpoint, err = newV1Endpoint(s.flags, allMiddlewares, s.websocketUpgrader, allModels)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		group.Path(
			v1Endpoint.Method(),
			v1Endpoint.Path(),
			v1Endpoint.Endpoint(),
		).Middlewares(v1Endpoint.Middlewares())

		group.Path(
			v1Endpoint.Users.Login.Method(),
			v1Endpoint.Users.Login.Path(),
			v1Endpoint.Users.Login.Endpoint(),
		).Middlewares(v1Endpoint.Users.Login.Middlewares())

		group.Path(
			v1Endpoint.Events.Lister.Method(),
			v1Endpoint.Events.Lister.Path(),
			v1Endpoint.Events.Lister.Endpoint(),
		).Middlewares(v1Endpoint.Events.Lister.Middlewares())
		group.Path(
			v1Endpoint.Events.Searcher.Method(),
			v1Endpoint.Events.Searcher.Path(),
			v1Endpoint.Events.Searcher.Endpoint(),
		).Middlewares(v1Endpoint.Events.Searcher.Middlewares())
		group.Path(
			v1Endpoint.Events.Watcher.Method(),
			v1Endpoint.Events.Watcher.Path(),
			v1Endpoint.Events.Watcher.Endpoint(),
		).Middlewares(v1Endpoint.Events.Watcher.Middlewares())
	}

	return s, nil
}

func (s *Server) Boot() error {
	if err := s.atreugo.ListenAndServe(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func newRootEndpoint(flags *flags.Flags, middleware *middleware.Middleware, models *models.Model) (*root.Endpoint, error) {
	var err error

	var endpoint *root.Endpoint
	{
		c := root.EndpointConfig{
			Flags:      flags,
			Middleware: middleware,
			Models:     models,
		}
		endpoint, err = root.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func newV1Endpoint(flags *flags.Flags, middleware *middleware.Middleware, websocketUpgrader *websocket.Upgrader, models *models.Model) (*v1.Endpoint, error) {
	var err error

	var endpoint *v1.Endpoint
	{
		c := v1.EndpointConfig{
			Flags:             flags,
			Middleware:        middleware,
			Models:            models,
			WebsocketUpgrader: websocketUpgrader,
		}
		endpoint, err = v1.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
