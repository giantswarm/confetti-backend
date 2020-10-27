package events

type ResponseEvent struct {
	Active bool `json:"active"`
	ID string `json:"id"`
	Name string `json:"name"`
	EventType string `json:"event_type"`
}

type Response struct {
	Events []ResponseEvent `json:"events"`
}