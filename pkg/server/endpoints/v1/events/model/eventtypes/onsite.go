package eventtypes

import "github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/model"

const (
	eventType = "onsite"
)

type OnsiteEventRoomID string

type OnsiteEventRoom struct {
	ID OnsiteEventRoomID
	Name string
	ConferenceURL string
}

type OnsiteEvent struct {
	*model.BaseEvent
	Rooms []OnsiteEventRoom
}

func (oe *OnsiteEvent) Type() model.EventType {
	return eventType
}

func NewOnsiteEvent() *OnsiteEvent {
	e := &OnsiteEvent{}

	return e
}