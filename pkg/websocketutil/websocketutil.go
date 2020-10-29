package websocketutil

import (
	"github.com/atreugo/websocket"
	"github.com/giantswarm/microerror"
)

func HandleConnection(connection *websocket.Conn, hub *Hub) error {
	client, err := NewClient(hub, connection)
	if err != nil {
		return microerror.Mask(err)
	}

	go client.WritePump()
	client.ReadPump()

	return nil
}
