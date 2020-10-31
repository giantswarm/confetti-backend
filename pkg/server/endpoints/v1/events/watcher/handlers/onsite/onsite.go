package onsite

import (
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/handlers"
	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/payloads"
	onsitePayload "github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/payloads/onsite"
	"github.com/giantswarm/confetti-backend/pkg/server/models"
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
	"github.com/giantswarm/confetti-backend/pkg/server/models/users/types"
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

	oeh.handleInitialStateMessages(event, message)
}

func (oeh *OnsiteEventHandler) OnClientDisconnect(message handlers.EventHandlerMessage) {
	event, err := oeh.findEventByID(message.EventID)
	if IsInvalidEventType(err) {
		return
	} else if err != nil {
		// TODO(axbarsan): Dispatch error message.
		return
	}

	oeh.handleFiniteStateMessages(event, message)
}

func (oeh *OnsiteEventHandler) OnClientMessage(message handlers.EventHandlerMessage) {
	event, err := oeh.findEventByID(message.EventID)
	if IsInvalidEventType(err) {
		return
	} else if err != nil {
		// TODO(axbarsan): Dispatch error message.
		return
	}

	payload := payloads.MessagePayload{}
	err = payload.Deserialize(message.Message.Payload)
	if err != nil {
		fmt.Println(err)
		// TODO(axbarsan): Dispatch error message.
	}

	switch payload.MessageType {
	case onsitePayload.OnsiteRoomJoinRequest:
		oeh.handleRoomJoin(event, message, payload)
	case onsitePayload.OnsiteRoomLeaveRequest:
		oeh.handleRoomLeave(event, message, payload)
	}
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

func (oek *OnsiteEventHandler) handleInitialStateMessages(event *eventsModelTypes.OnsiteEvent, message handlers.EventHandlerMessage) {
	// Join lobby.
	event.Lobby[message.User] = true

	var payloadBytes []byte
	var payload payloads.MessagePayload

	for _, room := range event.Rooms {
		payload = roomMessagePayload(
			onsitePayload.OnsiteRoomUpdateAttendeeCounter,
			"",
			room.ID,
			toIntPtr(len(room.Attendees)),
		)
		payloadBytes, _ = payload.Serialize()
		message.Message.Client.Emit(payloadBytes)
	}
}

func (oek *OnsiteEventHandler) handleFiniteStateMessages(event *eventsModelTypes.OnsiteEvent, message handlers.EventHandlerMessage) {
	// Leave lobby.
	delete(event.Lobby, message.User)

	var err error

	var payloadBytes []byte
	var payload payloads.MessagePayload

	var roomIndex int
	var room eventsModelTypes.OnsiteEventRoom
	for _, room = range event.Rooms {
		if _, ok := room.Attendees[message.User]; !ok {
			// User is not in the room.
			continue
		}

		delete(room.Attendees, message.User)
	}

	{
		event.Rooms[roomIndex] = room
		_, err = oek.models.Events.Update(event)
		if err != nil {
			// Ignore error, we don't want to send it to the client.
			return
		}
	}

	payloadBytes, _ = payload.Serialize()
	message.Message.Client.Emit(payloadBytes)

	// Broadcast room attendee counter update message.
	payload = roomMessagePayload(
		onsitePayload.OnsiteRoomUpdateAttendeeCounter,
		"",
		room.ID,
		toIntPtr(len(room.Attendees)),
	)
	payloadBytes, _ = payload.Serialize()
	message.Hub.BroadcastAll(payloadBytes)
}

func (oek *OnsiteEventHandler) handleRoomJoin(event *eventsModelTypes.OnsiteEvent, message handlers.EventHandlerMessage, messagePayload payloads.MessagePayload) {
	var success bool
	var err error

	var payloadBytes []byte
	var payload payloads.MessagePayload

	var roomIndex int
	var room eventsModelTypes.OnsiteEventRoom
	for roomIndex, room = range event.Rooms {
		if room.ID != messagePayload.Data.RoomID {
			continue
		}

		{
			if room.Attendees == nil {
				room.Attendees = make(map[*types.User]bool)
			}
			if _, ok := room.Attendees[message.User]; !ok {
				room.Attendees[message.User] = true
				delete(event.Lobby, message.User)
			} else {
				// User already in room.
				return
			}
		}

		payload = roomMessagePayload(
			onsitePayload.OnsiteRoomJoinSuccess,
			fmt.Sprintf("Joined room with ID '%s' successfully.", room.ID),
			room.ID,
			nil,
		)

		success = true

		break
	}

	if !success {
		payload = roomMessagePayload(
			onsitePayload.OnsiteRoomJoinError,
			fmt.Sprintf("Room with ID '%s' doesn't exist.", messagePayload.Data.RoomID),
			messagePayload.Data.RoomID,
			nil,
		)

		payloadBytes, _ = payload.Serialize()
		message.Message.Client.Emit(payloadBytes)

		return
	}

	{
		event.Rooms[roomIndex] = room

		_, err = oek.models.Events.Update(event)
		if err != nil {
			payload = roomMessagePayload(
				onsitePayload.OnsiteRoomJoinError,
				fmt.Sprintf("Couldn't join room with ID '%s'.", room.ID),
				room.ID,
				nil,
			)
		}
	}

	payloadBytes, _ = payload.Serialize()
	message.Message.Client.Emit(payloadBytes)

	// Broadcast room attendee counter update message.
	if success {
		payload = roomMessagePayload(
			onsitePayload.OnsiteRoomUpdateAttendeeCounter,
			"",
			room.ID,
			toIntPtr(len(room.Attendees)),
		)
		payloadBytes, _ = payload.Serialize()
		message.Hub.BroadcastAll(payloadBytes)
	}
}

func (oek *OnsiteEventHandler) handleRoomLeave(event *eventsModelTypes.OnsiteEvent, message handlers.EventHandlerMessage, messagePayload payloads.MessagePayload) {
	var success bool
	var err error

	var payloadBytes []byte
	var payload payloads.MessagePayload

	var roomIndex int
	var room eventsModelTypes.OnsiteEventRoom
	for roomIndex, room = range event.Rooms {
		if room.ID != messagePayload.Data.RoomID {
			continue
		}

		{
			if _, ok := room.Attendees[message.User]; ok && room.Attendees != nil {
				delete(room.Attendees, message.User)
			} else {
				// User not in room.
				return
			}
			event.Lobby[message.User] = true
		}

		payload = roomMessagePayload(
			onsitePayload.OnsiteRoomLeaveSuccess,
			fmt.Sprintf("Left room with ID '%s' successfully.", room.ID),
			room.ID,
			nil,
		)

		success = true

		break
	}

	if !success {
		payload = roomMessagePayload(
			onsitePayload.OnsiteRoomLeaveError,
			fmt.Sprintf("Room with ID '%s' doesn't exist.", messagePayload.Data.RoomID),
			messagePayload.Data.RoomID,
			nil,
		)

		payloadBytes, _ = payload.Serialize()
		message.Message.Client.Emit(payloadBytes)

		return
	}

	{
		event.Rooms[roomIndex] = room
		_, err = oek.models.Events.Update(event)
		if err != nil {
			payload = roomMessagePayload(
				onsitePayload.OnsiteRoomLeaveError,
				fmt.Sprintf("Couldn't leave room with ID '%s'.", room.ID),
				room.ID,
				nil,
			)
		}
	}

	payloadBytes, _ = payload.Serialize()
	message.Message.Client.Emit(payloadBytes)

	// Broadcast room attendee counter update message.
	if success {
		payload = roomMessagePayload(
			onsitePayload.OnsiteRoomUpdateAttendeeCounter,
			"",
			room.ID,
			toIntPtr(len(room.Attendees)),
		)
		payloadBytes, _ = payload.Serialize()
		message.Hub.BroadcastAll(payloadBytes)
	}
}

func roomMessagePayload(msgType payloads.MessagePayloadType, message string, roomID string, numOfAttendees *int) payloads.MessagePayload {
	return payloads.MessagePayload{
		MessageType: msgType,
		Data: payloads.MessagePayloadData{
			Message: message,
			OnsitePayload: onsitePayload.OnsitePayload{
				RoomID:  roomID,
				Counter: numOfAttendees,
			},
		},
	}
}

func toIntPtr(d int) *int {
	return &d
}
