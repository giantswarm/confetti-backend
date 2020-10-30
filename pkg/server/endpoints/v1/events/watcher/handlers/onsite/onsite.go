package onsite

import (
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type OnsiteEventConfig struct {
	Models *models.Model
}

type OnsiteEventHandler struct {
	models *models.Model
}

func NewOnsiteEvent(c OnsiteEventConfig) *OnsiteEventHandler {
	oeh := &OnsiteEventHandler{
		models: c.Models,
	}

	return oeh
}

func (oeh *OnsiteEventHandler) OnClientConnect(message handlers.EventHandlerMessage) {
	event, err := oeh.findEventByID(message.EventID)
	if IsInvalidEventType(err) {
		return
	} else if err != nil {
		// TODO(axbarsan): Dispatch error message.
		return
	}

	event.Lobby[message.User] = true

	// TODO(axbarsan): Dispatch success message.
}

func (oeh *OnsiteEventHandler) OnClientDisconnect(message handlers.EventHandlerMessage) {
	event, err := oeh.findEventByID(message.EventID)
	if IsInvalidEventType(err) {
		return
	} else if err != nil {
		// TODO(axbarsan): Dispatch error message.
		return
	}

	delete(event.Lobby, message.User)

	// TODO(axbarsan): Dispatch disconnect message.
}

func (oeh *OnsiteEventHandler) OnClientMessage(message handlers.EventHandlerMessage) {
	// TODO(axbarsan): Parse message into custom format.
	fmt.Println(string(message.Message.Payload))
}

func (oeh *OnsiteEventHandler) findEventByID(id string) (*eventsModelTypes.OnsiteEvent, error) {
	event, err := oeh.models.Events.FindOneByID(id)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	onsiteEvent, ok := event.(*eventsModelTypes.OnsiteEvent)
	if !ok {
		return nil, microerror.Mask(invalidEventTypeError)
	}

	return onsiteEvent, nil
}
