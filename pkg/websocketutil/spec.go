package websocketutil

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
//
// Taken from https://github.com/fasthttp/websocket/blob/master/_examples/chat/fasthttp/hub.go
// and modified to match our need.
type Hub interface {
	// Run starts the hub.
	Run()

	// On adds an event listener for a certain client-specific
	// event.
	On(event Event, callback EventCallback)
	// BroadcastAll sends a message to all connected clients.
	BroadcastAll(message []byte)
	// BroadcastAllExcept will send a message to all connected
	// clients, except the one specified.
	BroadcastAllExcept(message []byte, c *Client)

	// SendMessage transmits a message from the client
	// to the hub.
	SendMessage(message ClientMessage)

	// RegisterClient adds a new client to the hub.
	RegisterClient(client *Client)
	// UnregisterClient removes a known client from the hub.
	UnregisterClient(client *Client)
}
