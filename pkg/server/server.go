package server

import (
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/root"
	"github.com/giantswarm/confetti-backend/pkg/server/v1/users"
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

	var rootEndpoint *root.Endpoint
	{
		rootEndpoint, err = newRootEndpoint(c.Flags)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		s.atreugo.Path(rootEndpoint.Method(), rootEndpoint.Path(), rootEndpoint.Endpoint())
	}

	v1 := s.atreugo.NewGroupPath("/v1")

	var usersEndpoint *users.Endpoint
	{
		usersEndpoint, err = newUsersEndpoint(c.Flags)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		v1.Path(usersEndpoint.Method(), usersEndpoint.Path(), usersEndpoint.Endpoint())
		v1.Path(usersEndpoint.Login.Method(), usersEndpoint.Login.Path(), usersEndpoint.Login.Endpoint())
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

func newUsersEndpoint(flags *flags.Flags) (*users.Endpoint, error) {
	var err error

	var endpoint *users.Endpoint
	{
		c := users.EndpointConfig{
			Flags: flags,
		}
		endpoint, err = users.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
