package watcher

import (
	"github.com/atreugo/websocket"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/internal/flags"
	"github.com/giantswarm/confetti-backend/pkg/server/context/event"
	"github.com/giantswarm/confetti-backend/pkg/server/context/user"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers"
	eventHandlers "github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers/event"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
	usersModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/users/types"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
)

type ServiceConfig struct {
	Flags  *flags.Flags
	Hub    websocketutil.Hub
	Models *models.Model
}

type Service struct {
	flags                  *flags.Flags
	hub                    websocketutil.Hub
	models                 *models.Model
	eventHandlerCollection *handlers.EventHandlerCollection
}

func NewService(c ServiceConfig) (*Service, error) {
	if c.Flags == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flags must not be empty", c)
	}
	if c.Hub == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Hub must not be empty", c)
	}
	if c.Models == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Models must not be empty", c)
	}

	var ehc *handlers.EventHandlerCollection
	{
		defaultEventHandlerConfig := eventHandlers.DefaultEventConfig{
			Models: c.Models,
		}
		defaultEventHandler := eventHandlers.NewDefaultEventHandler(defaultEventHandlerConfig)

		onsiteHandlerConfig := eventHandlers.OnsiteEventConfig{
			Models: c.Models,
		}
		onsiteHandler := eventHandlers.NewOnsiteEventHandler(onsiteHandlerConfig)

		ehc = handlers.NewEventHandlerCollection()
		ehc.RegisterHandler((&eventsModelTypes.BaseEvent{}).Type(), defaultEventHandler)
		ehc.RegisterHandler(eventsModelTypes.NewOnsiteEvent().Type(), onsiteHandler)
	}

	service := &Service{
		flags:                  c.Flags,
		models:                 c.Models,
		hub:                    c.Hub,
		eventHandlerCollection: ehc,
	}

	{
		// Bind hub message handlers.
		c.Hub.On(websocketutil.EventConnected, service.handleClientConnect)
		c.Hub.On(websocketutil.EventDisconnected, service.handleClientDisconnect)
		c.Hub.On(websocketutil.EventMessage, service.handleClientMessage)

		go c.Hub.Run()
	}

	return service, nil
}

func (s *Service) NewClient(ws *websocket.Conn) error {
	c := websocketutil.ClientConfig{
		Hub:        s.hub,
		Connection: ws,
	}

	_, err := websocketutil.NewClient(c)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// handleClientConnect is executed whenever a new websocket
// connection is established.
func (s *Service) handleClientConnect(message websocketutil.ClientMessage) {
	handlerMessage := handlers.EventHandlerMessage{
		ClientMessage: message,
		Event:         s.getEvent(message),
		User:          s.getUser(message),
		Hub:           s.hub,
	}
	s.eventHandlerCollection.Visit(handlerMessage.Event.Type(), func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientConnect(handlerMessage)
	})
}

// handleClientDisconnect is executed whenever a websocket connection
// is about to be closed.
func (s *Service) handleClientDisconnect(message websocketutil.ClientMessage) {
	handlerMessage := handlers.EventHandlerMessage{
		ClientMessage: message,
		Event:         s.getEvent(message),
		User:          s.getUser(message),
		Hub:           s.hub,
	}
	s.eventHandlerCollection.Visit(handlerMessage.Event.Type(), func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientDisconnect(handlerMessage)
	})
}

// handleClientMessage is executed whenever a websocket client
// sends a message to the server.
func (s *Service) handleClientMessage(message websocketutil.ClientMessage) {
	handlerMessage := handlers.EventHandlerMessage{
		ClientMessage: message,
		Event:         s.getEvent(message),
		User:          s.getUser(message),
		Hub:           s.hub,
	}
	s.eventHandlerCollection.Visit(handlerMessage.Event.Type(), func(eventHandler handlers.EventHandler) {
		eventHandler.OnClientMessage(handlerMessage)
	})
}

func (s *Service) getEvent(message websocketutil.ClientMessage) eventsModelTypes.Event {
	// Event validation is already done in middleware.
	event, _ := event.FromValueGetter(message.Client.GetUserValue)

	return event
}

func (s *Service) getUser(message websocketutil.ClientMessage) *usersModelTypes.User {
	// User validation already done in middleware.
	user, _ := user.FromValueGetter(message.Client.GetUserValue)

	return user
}
