package mbcv

type commandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Status string
}
