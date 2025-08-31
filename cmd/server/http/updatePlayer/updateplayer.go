package updateplayer

import "darts-counter/models"

type Request struct {
	ID   string
	Name string
}

type Response struct {
	*models.Player
}
