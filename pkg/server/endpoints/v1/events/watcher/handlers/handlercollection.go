package handlers

import (
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

type EventHandlerCollection struct {
	handlers map[eventsModelTypes.EventType]EventHandler
}

func NewEventHandlerCollection() *EventHandlerCollection {
	ehc := &EventHandlerCollection{
		handlers: make(map[eventsModelTypes.EventType]EventHandler),
	}

	return ehc
}

// RegisterHandler registers a new event handler for a known
// event type.
//
// A event type can only have a single handler.
func (ehc *EventHandlerCollection) RegisterHandler(eventType eventsModelTypes.EventType, handler EventHandler) {
	ehc.handlers[eventType] = handler
}

// UnregisterHandler removes an event handler for an event type.
func (ehc *EventHandlerCollection) UnregisterHandler(eventType eventsModelTypes.EventType) {
	delete(ehc.handlers, eventType)
}

// Visit is using the visitor pattern to traverse all the
// registered event handlers.
func (ehc *EventHandlerCollection) Visit(visitor func(handler EventHandler)) {
	for _, handler := range ehc.handlers {
		visitor(handler)
	}
}
