package createplayer

import "darts-counter/models"

type Request struct {
	Name string
}

type Response struct {
	*models.Player
}
