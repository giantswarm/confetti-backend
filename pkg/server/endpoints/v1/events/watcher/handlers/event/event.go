package event

import (
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/payloads"
	eventPayloads "github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/payloads/event"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type DefaultEventConfig struct {
	Models *models.Model
	Logger micrologger.Logger
}

type DefaultEventHandler struct {
	models *models.Model
	logger micrologger.Logger
}

func NewDefaultEventHandler(c DefaultEventConfig) *DefaultEventHandler {
	oeh := &DefaultEventHandler{
		models: c.Models,
		logger: c.Logger,
	}

	return oeh
}

func (oeh *DefaultEventHandler) OnClientConnect(message handlers.EventHandlerMessage) {
	payload := getEventConfigurationMessagePayload(message.Event)
	payloadBytes, _ := payload.Serialize()
	message.ClientMessage.Client.Emit(payloadBytes)
}

func (oeh *DefaultEventHandler) OnClientDisconnect(message handlers.EventHandlerMessage) {
}

func (oeh *DefaultEventHandler) OnClientMessage(message handlers.EventHandlerMessage) {
}

func getEventConfigurationMessagePayload(event eventsModelTypes.Event) payloads.MessagePayload {
	var eventConfiguration *eventPayloads.Configuration
	{
		eventConfiguration = &eventPayloads.Configuration{
			Active:    event.Active(),
			ID:        event.ID(),
			Name:      event.Name(),
			EventType: string(event.Type()),
		}

		// Add event type-specific details.
		switch e := event.(type) {
		case *eventsModelTypes.OnsiteEvent:
			eventConfiguration.Details.Rooms = make([]eventPayloads.ConfigurationOnsiteRoom, 0, len(e.Rooms))
			for _, room := range e.Rooms {
				eventConfiguration.Details.Rooms = append(eventConfiguration.Details.Rooms, eventPayloads.ConfigurationOnsiteRoom{
					ID:            room.ID,
					Name:          room.Name,
					Description:   room.Description,
					ConferenceURL: room.ConferenceURL,
				})
			}
		}
	}

	payload := payloads.MessagePayload{
		MessageType: eventPayloads.EventUpdateConfiguration,
		Data: payloads.MessagePayloadData{
			Message: "Event configuration updated.",
			DefaultEventPayload: eventPayloads.DefaultEventPayload{
				Configuration: eventConfiguration,
			},
		},
	}

	return payload
}
