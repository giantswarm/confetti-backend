package websocketutil

const (
	// EventConnected will be executed when a client is connected.
	EventConnected = iota
	// EventDisconnected will be executed when a client is disconnected.
	EventDisconnected
	// EventMessage will be executed when a client sends a message to the server.
	EventMessage
)

type Event int

// EventCallback is a function that gets executed when a event is triggered.
type EventCallback func(clientMessage ClientMessage)
