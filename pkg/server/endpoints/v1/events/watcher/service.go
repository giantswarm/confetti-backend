package watcher

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model/types"
)

type ServiceConfig struct {
	Flags      *flags.Flags
	Repository *model.Repository
}

type Service struct {
	flags      *flags.Flags
	repository *model.Repository
}

func NewService(c ServiceConfig) (*Service, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Repository == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Repository must not be empty", c)
	}

	service := &Service{
		flags:      c.Flags,
		repository: c.Repository,
	}

	return service, nil
}

func (s *Service) GetEventByID(id string) (types.Event, error) {
	event, err := s.repository.FindOneByID(id)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return event, nil
}
