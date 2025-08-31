package creatematch

type Request struct {
	Pids      []string
	StartAt   int
	StartMode uint8
	EndMode   uint8
}
