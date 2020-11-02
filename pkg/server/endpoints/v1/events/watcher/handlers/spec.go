package handlers

import (
	usersModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/users/types"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
)

type EventHandlerMessage struct {
	websocketutil.ClientMessage

	EventID string
	User    *usersModelTypes.User
	Hub     *websocketutil.Hub
}

// EventHandler is the place where event type specific
// business logic would be written.
type EventHandler interface {
	// OnClientConnect runs when a client is connected.
	OnClientConnect(message EventHandlerMessage)
	// OnclientDisconnect runs right before a client is disconnected.
	OnClientDisconnect(message EventHandlerMessage)
	// OnClientMessage runs when the client sends a message to the server.
	OnClientMessage(message EventHandlerMessage)
}
