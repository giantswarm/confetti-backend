package events

import (
	usersModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/users/types"
)

const (
	eventType = "onsite"
)

type OnsiteEventRoom struct {
	ID            string
	Name          string
	Description   string
	ConferenceURL string
	Attendees     map[*usersModelTypes.User]bool
}

type OnsiteEvent struct {
	*BaseEvent

	Lobby map[*usersModelTypes.User]bool
	Rooms []OnsiteEventRoom
}

func (oe *OnsiteEvent) Type() EventType {
	return eventType
}

func NewOnsiteEvent() *OnsiteEvent {
	be := &BaseEvent{}

	e := &OnsiteEvent{
		BaseEvent: be,
		Lobby:     make(map[*usersModelTypes.User]bool),
	}

	return e
}
