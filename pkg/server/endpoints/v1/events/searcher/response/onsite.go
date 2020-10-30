package response

type ResponseOnsiteRoom struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ConferenceURL string `json:"conference_url"`
}
