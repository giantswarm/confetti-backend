package watcher

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/context/user"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers/onsite"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	events "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
	usersModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/users/types"
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
		onsiteHandlerConfig := onsite.OnsiteEventConfig{
			Models: c.Models,
		}
		onsiteHandler := onsite.NewOnsiteEvent(onsiteHandlerConfig)

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
	handlerMessage := handlers.EventHandlerMessage{
		EventID: s.getEventID(message),
		User:    s.getUser(message),
		Message: message,
	}
	s.eventHandlerCollection.Visit(func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientConnect(handlerMessage)
	})
}

func (s *Service) HandleClientDisconnect(message websocketutil.ClientMessage) {
	handlerMessage := handlers.EventHandlerMessage{
		EventID: s.getEventID(message),
		User:    s.getUser(message),
		Message: message,
	}
	s.eventHandlerCollection.Visit(func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientDisconnect(handlerMessage)
	})
}

func (s *Service) HandleClientMessage(message websocketutil.ClientMessage) {
	handlerMessage := handlers.EventHandlerMessage{
		EventID: s.getEventID(message),
		User:    s.getUser(message),
		Message: message,
	}
	s.eventHandlerCollection.Visit(func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientMessage(handlerMessage)
	})
}

func (s *Service) getEventID(message websocketutil.ClientMessage) string {
	// ID validation is already done in middleware.
	return message.Client.GetUserValue("id").(string)
}

func (s *Service) getUser(message websocketutil.ClientMessage) *usersModelTypes.User {
	// User validation already done in middleware.
	user, _ := user.FromUserValueGetter(message.Client.GetUserValue)

	return user
}
