package root

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/flags"
	"github.com/giantswarm/confetti-backend/pkg/project"
)

const (
	method = "GET"
	path   = "/"
)

type EndpointConfig struct {
	Flags *flags.Flags
}

type Endpoint struct {
	flags *flags.Flags
}

func NewEndpoint(c EndpointConfig) (*Endpoint, error) {
	endpoint := &Endpoint{
		flags: c.Flags,
	}

	return endpoint, nil
}

func (e *Endpoint) Endpoint() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		res := Response{
			Description: project.Description(),
			Name:        project.Name(),
			GitCommit:   project.GitSHA(),
			Source:      project.Source(),
			Version:     project.Version(),
		}

		err := ctx.JSONResponse(res, http.StatusOK)
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