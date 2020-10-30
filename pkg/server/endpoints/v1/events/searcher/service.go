package searcher

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	eventsModel "github.com/giantswarm/confetti-backend/pkg/server/models/events"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type ServiceConfig struct {
	Flags      *flags.Flags
	Repository *eventsModel.Repository
}

type Service struct {
	flags      *flags.Flags
	repository *eventsModel.Repository
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

func (s *Service) GetEventByID(id string) (eventsModelTypes.Event, error) {
	event, err := s.repository.FindOneByID(id)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return event, nil
}
