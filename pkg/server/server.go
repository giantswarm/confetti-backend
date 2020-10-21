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

func New(c Config) *Server {
	s := &Server{
		atreugo: c.Atreugo,
		flags:   c.Flags,
	}

	return s
}

func (s *Server) Boot() error {
	var err error

	// v1 := s.atreugo.NewGroupPath("/v1")

	var rootEndpoint *root.Endpoint
	{
		rootEndpoint, err = s.newRootEndpoint()
		if err != nil {
			return microerror.Mask(err)
		}

		s.atreugo.Path(rootEndpoint.Method(), rootEndpoint.Path(), rootEndpoint.Endpoint())
	}

	if err = s.atreugo.ListenAndServe(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (s *Server) newRootEndpoint() (*root.Endpoint, error) {
	c := root.EndpointConfig{
		Flags: s.flags,
	}
	e, err := root.NewEndpoint(&c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return e, nil
}
