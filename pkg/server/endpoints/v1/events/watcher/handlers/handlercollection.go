package handlers

import (
	eventsModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/events/types"
)

var (
	// defaultEventType is the handler type that will be ran for
	// all event types.
	defaultEventType = (&eventsModelTypes.BaseEvent{}).Type()
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
// event handlers registered for a given event type.
func (ehc *EventHandlerCollection) Visit(eventType eventsModelTypes.EventType, visitor func(handler EventHandler)) {
	for handlerEventType, handler := range ehc.handlers {
		if handlerEventType != eventType && handlerEventType != defaultEventType {
			continue
		}

		visitor(handler)
	}
}
