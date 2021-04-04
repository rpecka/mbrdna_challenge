package server

var intentMap = map[string]Command{
	"CAR.unlock": unlockCar,
	"CAR.lock":   lockCar,
	"CAR.find":   findCar,
}

type IntentHandler interface {
	commandForIntent(intent string) *Command
}

type BasicIntentHandler struct{}

func (b BasicIntentHandler) commandForIntent(intent string) *Command {
	command, ok := intentMap[intent]
	if !ok {
		return nil
	}
	return &command
}
