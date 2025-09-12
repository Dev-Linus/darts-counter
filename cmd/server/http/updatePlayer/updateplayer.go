package updateplayer

import "darts-counter/models"

// Request represents an update player request payload.
type Request struct {
	ID   string
	Name string
}

// Response wraps the updated player entity returned to the client.
type Response struct {
	*models.Player
}
