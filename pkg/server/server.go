package server

import (
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/server/root"
	v1 "github.com/giantswarm/confetti-backend/pkg/server/v1"
)

type Config struct {
	Atreugo *atreugo.Atreugo
	Flags   *flags.Flags
}

type Server struct {
	atreugo *atreugo.Atreugo
	flags   *flags.Flags
}

func New(c Config) (*Server, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Atreugo == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Atreugo must not be empty", c)
	}

	var err error

	s := &Server{
		atreugo: c.Atreugo,
		flags:   c.Flags,
	}

	var allMiddlewares *middleware.Middleware
	{
		c := middleware.Config{
			Flags: s.flags,
		}
		allMiddlewares, err = middleware.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var rootEndpoint *root.Endpoint
	{
		rootEndpoint, err = newRootEndpoint(s.flags, allMiddlewares)
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

		v1Endpoint, err = newV1Endpoint(s.flags, allMiddlewares)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		group.Path(
			v1Endpoint.Method(),
			v1Endpoint.Path(),
			v1Endpoint.Endpoint(),
		).Middlewares(v1Endpoint.Middlewares())
		group.Path(
			v1Endpoint.Users.Method(),
			v1Endpoint.Users.Path(),
			v1Endpoint.Users.Endpoint(),
		).Middlewares(v1Endpoint.Users.Middlewares())
		group.Path(
			v1Endpoint.Users.Login.Method(),
			v1Endpoint.Users.Login.Path(),
			v1Endpoint.Users.Login.Endpoint(),
		).Middlewares(v1Endpoint.Users.Login.Middlewares())
	}

	return s, nil
}

func (s *Server) Boot() error {
	if err := s.atreugo.ListenAndServe(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func newRootEndpoint(flags *flags.Flags, middleware *middleware.Middleware) (*root.Endpoint, error) {
	var err error

	var endpoint *root.Endpoint
	{
		c := root.EndpointConfig{
			Flags:      flags,
			Middleware: middleware,
		}
		endpoint, err = root.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func newV1Endpoint(flags *flags.Flags, middleware *middleware.Middleware) (*v1.Endpoint, error) {
	var err error

	var endpoint *v1.Endpoint
	{
		c := v1.EndpointConfig{
			Flags:      flags,
			Middleware: middleware,
		}
		endpoint, err = v1.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
