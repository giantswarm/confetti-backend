package v1

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/v1/users"
)

const (
	method = "GET"
	path   = "/"
)

type EndpointConfig struct {
	Flags *flags.Flags
}

type Endpoint struct {
	Users *users.Endpoint

	flags *flags.Flags
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	usersEndpoint, err := createUsersEndpoint(c.Flags)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	endpoint := &Endpoint{
		Users: usersEndpoint,

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

func createUsersEndpoint(flags *flags.Flags) (*users.Endpoint, error) {
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
