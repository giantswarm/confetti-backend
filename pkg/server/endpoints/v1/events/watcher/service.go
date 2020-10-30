package watcher

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
)

type ServiceConfig struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Service struct {
	flags  *flags.Flags
	models *models.Model
}

func NewService(c ServiceConfig) (*Service, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}

	service := &Service{
		flags:  c.Flags,
		models: c.Models,
	}

	return service, nil
}
