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

func (ehc *EventHandlerCollection) RegisterHandler(eventType eventsModelTypes.EventType, handler EventHandler) {
	ehc.handlers[eventType] = handler
}

func (ehc *EventHandlerCollection) UnregisterHandler(eventType eventsModelTypes.EventType) {
	delete(ehc.handlers, eventType)
}

func (ehc *EventHandlerCollection) Visit(visitor func(handler EventHandler)) {
	for _, handler := range ehc.handlers {
		visitor(handler)
	}
}
