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
		rootEndpoint, err = newRootEndpoint(s.flags)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		s.atreugo.Path(rootEndpoint.Method(), rootEndpoint.Path(), rootEndpoint.Endpoint())
	}

	var v1Endpoint *v1.Endpoint
	{
		group := s.atreugo.NewGroupPath("/v1")

		v1Endpoint, err = newV1Endpoint(s.flags)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		group.Path(v1Endpoint.Method(), v1Endpoint.Path(), v1Endpoint.Endpoint())
		group.Path(v1Endpoint.Users.Method(), v1Endpoint.Users.Path(), v1Endpoint.Users.Endpoint())
		group.Path(
			v1Endpoint.Users.Login.Method(),
			v1Endpoint.Users.Login.Path(),
			v1Endpoint.Users.Login.Endpoint(),
		).UseBefore(allMiddlewares.Authentication.Middleware)
	}

	return s, nil
}

func (s *Server) Boot() error {
	if err := s.atreugo.ListenAndServe(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func newRootEndpoint(flags *flags.Flags) (*root.Endpoint, error) {
	var err error

	var endpoint *root.Endpoint
	{
		c := root.EndpointConfig{
			Flags: flags,
		}
		endpoint, err = root.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}

func newV1Endpoint(flags *flags.Flags) (*v1.Endpoint, error) {
	var err error

	var endpoint *v1.Endpoint
	{
		c := v1.EndpointConfig{
			Flags: flags,
		}
		endpoint, err = v1.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
