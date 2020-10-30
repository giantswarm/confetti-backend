package watcher

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	events "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
)

type ServiceConfig struct {
	Flags  *flags.Flags
	Models *models.Model
}

type Service struct {
	flags                  *flags.Flags
	models                 *models.Model
	eventHandlerCollection *handlers.EventHandlerCollection
}

func NewService(c ServiceConfig) (*Service, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}

	var ehc *handlers.EventHandlerCollection
	{
		onsiteHandlerConfig := handlers.OnsiteEventConfig{
			Models: c.Models,
		}
		onsiteHandler := handlers.NewOnsiteEvent(onsiteHandlerConfig)

		ehc = handlers.NewEventHandlerCollection()
		ehc.RegisterHandler(events.NewOnsiteEvent().Type(), onsiteHandler)
	}

	service := &Service{
		flags:                  c.Flags,
		models:                 c.Models,
		eventHandlerCollection: ehc,
	}

	return service, nil
}

func (s *Service) HandleClientConnect(message websocketutil.ClientMessage) {
	id := s.getEventID(message)
	s.eventHandlerCollection.Visit(func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientConnect(id, message)
	})
}

func (s *Service) HandleClientDisconnect(message websocketutil.ClientMessage) {
	id := s.getEventID(message)
	s.eventHandlerCollection.Visit(func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientDisconnect(id, message)
	})
}

func (s *Service) HandleClientMessage(message websocketutil.ClientMessage) {
	id := s.getEventID(message)
	s.eventHandlerCollection.Visit(func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientMessage(id, message)
	})
}

func (s *Service) getEventID(message websocketutil.ClientMessage) string {
	// ID validation is already done in middleware.
	return message.Client.GetUserValue("id").(string)
}
