package response

type ResponseDetails struct {
	Rooms []ResponseOnsiteRoom `json:"rooms,omitempty"`
}

type Response struct {
	Active bool `json:"active"`
	ID string `json:"id"`
	Name string `json:"name"`
	EventType string `json:"event_type"`
	Details ResponseDetails `json:"details"`
}