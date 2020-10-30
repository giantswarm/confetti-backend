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
	Attendees     []*usersModelTypes.User
}

type OnsiteEvent struct {
	*BaseEvent

	Lobby []*usersModelTypes.User
	Rooms []OnsiteEventRoom
}

func (oe *OnsiteEvent) Type() EventType {
	return eventType
}

func NewOnsiteEvent() *OnsiteEvent {
	be := &BaseEvent{}

	e := &OnsiteEvent{
		BaseEvent: be,
	}

	return e
}
