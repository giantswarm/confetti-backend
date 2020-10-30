package websocketutil

import (
	"github.com/giantswarm/microerror"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
//
// Originally from https://github.com/fasthttp/websocket/blob/master/_examples/chat/fasthttp/hub.go
// and modified to match our need.
type Hub struct {
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

func NewHub() (*Hub, error) {
	hc, err := NewHookCollection()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	h := &Hub{
		broadcast:      make(chan ClientMessage),
		register:       make(chan ClientMessage),
		unregister:     make(chan ClientMessage),
		clients:        make(map[*Client]bool),
		hookCollection: hc,
	}

	return h, nil
}

func (h *Hub) Run() {
	for {
		select {
		case clientMessage := <-h.register:
			h.addClient(clientMessage.Client)
			h.hookCollection.Call(EventConnected, clientMessage)
		case clientMessage := <-h.unregister:
			if _, ok := h.clients[clientMessage.Client]; ok {
				h.removeClient(clientMessage.Client)
			}
		case clientMessage := <-h.broadcast:
			h.hookCollection.Call(EventMessage, clientMessage)
		}
	}
}

func (h *Hub) On(event Event, callback EventCallback) {
	h.hookCollection.Register(event, callback)
}

func (h *Hub) BroadcastAll(message []byte) {
	for client := range h.clients {
		h.tryMessageClient(client, message)
	}
}

func (h *Hub) BroadcastAllExcept(message []byte, c *Client) {
	for client := range h.clients {
		if client == c {
			continue
		}

		h.tryMessageClient(client, message)
	}
}

func (h *Hub) addClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) removeClient(client *Client) {
	delete(h.clients, client)
	close(client.send)
}

func (h *Hub) tryMessageClient(client *Client, message []byte) {
	if isOpen := client.Emit(message); !isOpen {
		h.removeClient(client)
	}
}
