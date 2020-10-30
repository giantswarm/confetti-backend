package handlers

import (
	"fmt"

	"github.com/giantswarm/confetti-backend/pkg/server/models"
	"github.com/giantswarm/confetti-backend/pkg/websocketutil"
)

type OnsiteEventConfig struct {
	Models *models.Model
}

type OnsiteEventHandler struct {
	models *models.Model
}

func NewOnsiteEvent(c OnsiteEventConfig) *OnsiteEventHandler {
	oeh := &OnsiteEventHandler{
		models: c.Models,
	}

	return oeh
}

func (oeh *OnsiteEventHandler) OnClientConnect(eventID string, message websocketutil.ClientMessage) {
	fmt.Printf("connected to %s\n", eventID)
}

func (oeh *OnsiteEventHandler) OnClientDisconnect(eventID string, message websocketutil.ClientMessage) {
	fmt.Printf("disconnected from %s\n", eventID)
}

func (oeh *OnsiteEventHandler) OnClientMessage(eventID string, message websocketutil.ClientMessage) {
	fmt.Println(message.Payload)
}
