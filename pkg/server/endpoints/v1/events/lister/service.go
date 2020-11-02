package lister

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type ServiceConfig struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Service struct {
	flags  *flags.Flags
	Models *models.Model
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
		Models: c.Models,
	}

	return service, nil
}

func (s *Service) GetEvents() ([]eventsModelTypes.Event, error) {
	events, err := s.Models.Events.FindAll()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return events, nil
}
