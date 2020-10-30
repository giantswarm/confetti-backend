package payloads

import (
	"encoding/json"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/payloads/onsite"
)

type MessagePayloadType string

type MessagePayloadData struct {
	Message string `json:"message,omitempty"`

	onsite.OnsitePayload
}

type MessagePayload struct {
	MessageType MessagePayloadType `json:"type"`
	Data        MessagePayloadData `json:"data,omitempty"`
}

func (mp *MessagePayload) Serialize() ([]byte, error) {
	payloadBytes, err := json.Marshal(mp)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return payloadBytes, nil
}

func (mp *MessagePayload) Deserialize(payloadBytes []byte) error {
	err := json.Unmarshal(payloadBytes, mp)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
