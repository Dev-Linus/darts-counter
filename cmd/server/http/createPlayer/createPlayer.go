package createplayer

import "darts-counter/models"

// Request represents a create player request payload.
type Request struct {
	Name string
}

// Response wraps the created player entity returned to the client.
type Response struct {
	*models.Player
}
