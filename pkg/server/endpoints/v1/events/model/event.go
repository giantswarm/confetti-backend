package model

const (
	defaultEventType = "default"
)

type EventID string

type EventType string

type Event interface {
	Type() EventType
}

type BaseEvent struct {
	Active bool
	ID EventID
	Name string
}

func (be *BaseEvent) Type() EventType {
	return defaultEventType
}