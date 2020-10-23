package server

import (
	"github.com/giantswarm/microerror"
	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flag"
	"github.com/giantswarm/confetti-backend/pkg/server/root"
)

type Config struct {
	Atreugo *atreugo.Atreugo
	Flags   *flag.Flag
}

type Server struct {
	atreugo *atreugo.Atreugo
	flags   *flag.Flag
}

func New(c Config) (*Server, error) {
	var err error

	s := &Server{
		atreugo: c.Atreugo,
		flags:   c.Flags,
	}

	// v1 := s.atreugo.NewGroupPath("/v1")

	var rootEndpoint *root.Endpoint
	{
		rootEndpoint, err = newRootEndpoint(c.Flags)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		s.atreugo.Path(rootEndpoint.Method(), rootEndpoint.Path(), rootEndpoint.Endpoint())
	}

	return s, nil
}

func (s *Server) Boot() error {
	if err := s.atreugo.ListenAndServe(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func newRootEndpoint(flags *flag.Flag) (*root.Endpoint, error) {
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
