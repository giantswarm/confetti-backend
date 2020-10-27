package model

type EventID string

type EventType string

type Event struct {
	Active bool
	ID EventID
	EventType EventType
	Name string
}