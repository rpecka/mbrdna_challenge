package requests

type DoorCommand string

const (
	LOCK   = DoorCommand("LOCK")
	UNLOCK = DoorCommand("UNLOCK")
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Status string
}
