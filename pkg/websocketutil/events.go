package websocketutil

const (
	EventConnected = iota
	EventDisconnected
	EventMessage
)

type Event int

type EventCallback func(clientMessage ClientMessage)
