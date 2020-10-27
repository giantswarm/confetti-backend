package types

const (
	defaultEventType = "default"
)

type EventType string

type Event interface {
	Active() bool
	ID() string
	Name() string
	Type() EventType
}

type BaseEvent struct {
	active bool
	id     string
	name   string
}

func (be *BaseEvent) Active() bool {
	return be.active
}

func (be *BaseEvent) ID() string {
	return be.id
}

func (be *BaseEvent) Name() string {
	return be.name
}

func (be *BaseEvent) Type() EventType {
	return defaultEventType
}
