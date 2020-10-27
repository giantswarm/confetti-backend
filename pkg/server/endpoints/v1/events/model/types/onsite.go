package types

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
	*BaseEvent
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