package getmatch

import "darts-counter/models"

// Request represents a create player request payload.
type Request struct {
	Mid string
}

// Response wraps the created player entity returned to the client.
type Response struct {
	Match   *models.Match   `json:"match"`
	History *models.History `json:"history"`
}
