package login

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/microerror"
)

const (
	method = "POST"
	path   = "/users/login/"
)

type EndpointConfig struct {
	Flags   *flags.Flags
	Service *Service
}

type Endpoint struct {
	flags   *flags.Flags
	service *Service
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	endpoint := &Endpoint{
		flags:   c.Flags,
		service: c.Service,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		token, err := e.service.Authenticate()
		if err != nil {
			return microerror.Mask(err)
		}

		res := Response{
			Token: token,
		}

		err = ctx.JSONResponse(res, http.StatusOK)
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
