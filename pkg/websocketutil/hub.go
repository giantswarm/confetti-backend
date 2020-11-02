package websocketutil

import (
	"github.com/giantswarm/microerror"
)

type SocketHub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan ClientMessage

	// Register requests from the clients.
	register chan ClientMessage

	// Unregister requests from clients.
	unregister chan ClientMessage

	hookCollection *HookCollection
}

func NewSocketHub() (Hub, error) {
	hc, err := NewHookCollection()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	h := &SocketHub{
		broadcast:      make(chan ClientMessage),
		register:       make(chan ClientMessage),
		unregister:     make(chan ClientMessage),
		clients:        make(map[*Client]bool),
		hookCollection: hc,
	}

	return h, nil
}

func (sh *SocketHub) Run() {
	for {
		select {
		case clientMessage := <-sh.register:
			sh.addClient(clientMessage.Client)
			sh.hookCollection.Call(EventConnected, clientMessage)
		case clientMessage := <-sh.unregister:
			if _, ok := sh.clients[clientMessage.Client]; ok {
				sh.hookCollection.Call(EventDisconnected, clientMessage)
				sh.removeClient(clientMessage.Client)
			}
		case clientMessage := <-sh.broadcast:
			sh.hookCollection.Call(EventMessage, clientMessage)
		}
	}
}

func (sh *SocketHub) On(event Event, callback EventCallback) {
	sh.hookCollection.Register(event, callback)
}

func (sh *SocketHub) BroadcastAll(message []byte) {
	for client := range sh.clients {
		sh.tryMessageClient(client, message)
	}
}

func (sh *SocketHub) BroadcastAllExcept(message []byte, c *Client) {
	for client := range sh.clients {
		if client == c {
			continue
		}

		sh.tryMessageClient(client, message)
	}
}

func (sh *SocketHub) RegisterClient(client *Client) {
	clientMessage := ClientMessage{
		Client: client,
	}
	sh.register <- clientMessage
}

func (sh *SocketHub) UnregisterClient(client *Client) {
	clientMessage := ClientMessage{
		Client: client,
	}
	sh.unregister <- clientMessage
}

func (sh *SocketHub) SendMessage(message ClientMessage) {
	sh.broadcast <- message
}

func (sh *SocketHub) addClient(client *Client) {
	sh.clients[client] = true
}

func (sh *SocketHub) removeClient(client *Client) {
	delete(sh.clients, client)
	close(client.send)
}

// tryMessageClient will try to send a message to a
// specific client, if the connection is still open.
func (sh *SocketHub) tryMessageClient(client *Client, message []byte) {
	if isOpen := client.Emit(message); !isOpen {
		sh.removeClient(client)
	}
}
