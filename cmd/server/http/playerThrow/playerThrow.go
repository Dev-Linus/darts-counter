package playerthrow

import "darts-counter/models"

type Request struct {
	Pid    string
	Mid    string
	Amount uint32
}

type Response struct {
	Won            bool
	NextThrowBy    string
	Scores         map[string]uint32
	PossibleFinish []*models.Throws
}
