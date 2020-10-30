package handlers

import (
	usersModelTypes "github.com/giantswarm/confetti-backend/pkg/server/models/users/types"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
)

type EventHandlerMessage struct {
	EventID string
	User    *usersModelTypes.User
	Message websocketutil.ClientMessage
	Hub     *websocketutil.Hub
}

type EventHandler interface {
	OnClientConnect(message EventHandlerMessage)
	OnClientDisconnect(message EventHandlerMessage)
	OnClientMessage(message EventHandlerMessage)
}
