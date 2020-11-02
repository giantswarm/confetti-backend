package event

const (
	// EventUpdateConfiguration is broadcasted by
	// the server to all clients whenever something
	// changes in the configuration of an event.
	//
	// Also, this message is sent to each client,
	// when a new connection is established.
	EventUpdateConfiguration = "EVENT_UPDATE_CONFIGURATION"
)

type DefaultEventPayload struct {
	Configuration *Configuration `json:"configuration,omitempty"`
}
