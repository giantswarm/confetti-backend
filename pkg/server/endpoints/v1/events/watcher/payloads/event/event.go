package event

const (
	// EventUpdateConfiguration is broadcasted by
	// the server to all clients whenever something
	// changes in the configuration of an event.
	//
	// Also, this message is sent to each client,
	// when a new connection is established.
	EventUpdateConfiguration = "EVENT_UPDATE_CONFIGURATION"
	// EventInvalidPayloadError is sent by the server
	// when the client tries to send a message with
	// an invalid payload.
	EventInvalidPayloadError = "EVENT_INVALID_PAYLOAD_ERROR"
)

type DefaultEventPayload struct {
	Configuration *Configuration `json:"configuration,omitempty"`
}
