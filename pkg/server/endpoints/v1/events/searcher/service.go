package searcher

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type ServiceConfig struct {
	Flags  *flags.Flags
	Models *models.Model
	Logger micrologger.Logger
}

type Service struct {
	flags  *flags.Flags
	models *models.Model
	logger micrologger.Logger
}

func NewService(c ServiceConfig) (*Service, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", c)
	}

	service := &Service{
		flags:  c.Flags,
		models: c.Models,
		logger: c.Logger,
	}

	return service, nil
}

func (s *Service) GetEventByID(id string) (eventsModelTypes.Event, error) {
	event, err := s.models.Events.FindOneByID(id)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return event, nil
}
