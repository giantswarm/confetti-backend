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

func (oeh *OnsiteEventHandler) OnClientConnect(message websocketutil.ClientMessage) {
	fmt.Println("connected")
}

func (oeh *OnsiteEventHandler) OnClientDisconnect(message websocketutil.ClientMessage) {
	fmt.Println("disconnected")
}

func (oeh *OnsiteEventHandler) OnClientMessage(message websocketutil.ClientMessage) {
	fmt.Println(message.Payload)
}
