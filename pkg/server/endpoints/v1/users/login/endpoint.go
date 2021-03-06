package login

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/middleware"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

const (
	method = "POST"
	path   = "/users/login/"
)

type EndpointConfig struct {
	Flags      *flags.Flags
	Service    *Service
	Middleware *middleware.Middleware
	Models     *models.Model
	Logger     micrologger.Logger
}

type Endpoint struct {
	flags      *flags.Flags
	service    *Service
	middleware *middleware.Middleware
	models     *models.Model
	logger     micrologger.Logger
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Service == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Service must not be empty", c)
	}
	if c.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Middleware must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", c)
	}

	endpoint := &Endpoint{
		flags:      c.Flags,
		service:    c.Service,
		middleware: c.Middleware,
		models:     c.Models,
		logger:     c.Logger,
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

		return ctx.JSONResponse(res, http.StatusOK)
	}
}

func (e *Endpoint) Middlewares() atreugo.Middlewares {
	return atreugo.Middlewares{}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
