package event

type ConfigurationDetails struct {
	// Onsite specific.
	Rooms []ConfigurationOnsiteRoom `json:"rooms,omitempty"`
}

type Configuration struct {
	Active    bool                 `json:"active"`
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	EventType string               `json:"event_type"`
	Details   ConfigurationDetails `json:"details"`
}
