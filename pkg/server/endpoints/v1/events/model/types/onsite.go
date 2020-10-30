package types

const (
	eventType = "onsite"
)

type OnsiteEventRoom struct {
	ID            string
	Name          string
	Description   string
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
