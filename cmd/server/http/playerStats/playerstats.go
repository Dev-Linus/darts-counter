package playerstats

import "darts-counter/models"

type Request struct {
	Pid string
}

type Response struct {
	Name          string
	Throws        int
	Matches       int
	ActiveMatches int
	WinRate       float32
	MeanThrow     float32
	HighestFinish uint32
	Nemesis       *models.Player
	Dominating    *models.Player
}
