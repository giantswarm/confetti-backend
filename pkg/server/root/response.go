package root

type Response struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	GitCommit   string `json:"git_commit"`
	Source      string `json:"source"`
	Version     string `json:"version"`
}
