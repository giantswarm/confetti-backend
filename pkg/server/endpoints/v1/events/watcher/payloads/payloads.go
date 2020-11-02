package payloads

import (
	"encoding/json"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/confetti-backend/pkg/server/endpoints/v1/events/watcher/payloads/event"
)

type MessagePayloadType string

// MessagePayloadData contains all the keys
// that can be used in the payload of an
// event-specific websocket message.
type MessagePayloadData struct {
	event.EventPayload
	event.OnsitePayload

	// Message contains a user-friendly
	// message explaining the outcome
	// of the request.
	Message string `json:"message,omitempty"`
}

type MessagePayload struct {
	// MessageType represents the type of the message
	// sent, such as `EVENT_UPDATE_CONFIGURATION`.
	MessageType MessagePayloadType `json:"type"`
	// Data represents data passed inside a
	// message's payload.
	Data MessagePayloadData `json:"data,omitempty"`
}

// Serialize marshals the payload into a JSON-encoded
// byte buffer.
func (mp *MessagePayload) Serialize() ([]byte, error) {
	payloadBytes, err := json.Marshal(mp)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return payloadBytes, nil
}

// Deserialize unmarshals the JSON-encoded payload
// into the struct.
func (mp *MessagePayload) Deserialize(payloadBytes []byte) error {
	err := json.Unmarshal(payloadBytes, mp)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
