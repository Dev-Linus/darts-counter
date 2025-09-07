package playerthrow

import "darts-counter/models"

type Request struct {
	Pid   string
	Mid   string
	Throw models.ThrowType
}

type Response struct {
	Won            bool
	NextThrowBy    string
	Scores         map[string]int
	PossibleFinish []models.ThrowType
}
