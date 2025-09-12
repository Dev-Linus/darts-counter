package creatematch

// Request represents a create match request payload.
type Request struct {
	Pids      []string
	StartAt   int
	StartMode uint8
	EndMode   uint8
}
