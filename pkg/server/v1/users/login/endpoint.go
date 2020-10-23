package login

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flag"
)

const (
	method = "POST"
	path   = "/login/"
)

type EndpointConfig struct {
	Flags   *flag.Flag
	Service *Service
}

type Endpoint struct {
	flags   *flag.Flag
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
		res := Response{
			Token: "",
		}

		return ctx.JSONResponse(res, http.StatusOK)
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
