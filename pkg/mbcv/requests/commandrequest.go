package requests

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Status string
}
