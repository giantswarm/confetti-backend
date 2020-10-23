package users

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flag"
	"github.com/giantswarm/confetti-backend/pkg/server/v1/users/login"
	"github.com/giantswarm/microerror"
)

const (
	method = "GET"
	path   = "/users/"
)

type EndpointConfig struct {
	Flags *flag.Flag
}

type Endpoint struct {
	Login *login.Endpoint

	flags *flag.Flag
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	loginEndpoint, err := createLoginEndpoint(c.Flags)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Login: loginEndpoint,

		flags: c.Flags,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		err := ctx.HTTPResponse("", http.StatusNotFound)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}

func createLoginEndpoint(flags *flag.Flag) (*login.Endpoint, error) {
	var err error

	var service *login.Service
	{
		c := login.ServiceConfig{
			Flags: flags,
		}
		service, err = login.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var endpoint *login.Endpoint
	{
		c := login.EndpointConfig{
			Flags:   flags,
			Service: service,
		}
		endpoint, err = login.NewEndpoint(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return endpoint, nil
}
