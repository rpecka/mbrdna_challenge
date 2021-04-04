package server

type Command string

const (
	unlockCar = Command("UNLOCK")
	lockCar   = Command("LOCK")
	findCar   = Command("FIND")
)
