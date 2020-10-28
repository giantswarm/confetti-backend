package websocketutil

import (
	"github.com/atreugo/websocket"
	"github.com/giantswarm/microerror"
)

func HandleConnection(connection *websocket.Conn, hub *Hub) error {
	client, err := NewClient(hub, connection, make(chan []byte, 256))
	if err != nil {
		return microerror.Mask(err)
	}

	client.hub.register <- client

	go client.WritePump()
	client.ReadPump()

	return nil
}
