package root

import (
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"github.com/giantswarm/confetti-backend/flag"
	"github.com/giantswarm/confetti-backend/pkg/project"
)

const (
	method = "GET"
	path   = "/"
)

type EndpointConfig struct {
	Flags *flag.Flag
}

type Endpoint struct {
	flags *flag.Flag
}

func NewEndpoint(c *EndpointConfig) (*Endpoint, error) {
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

		return ctx.JSONResponse(res, http.StatusOK)
	}
}

func (e *Endpoint) Path() string {
	return path
}

func (e *Endpoint) Method() string {
	return method
}
