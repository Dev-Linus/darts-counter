package models

// MatchPlayer stores per-match stats for a single player.
type MatchPlayer struct {
	Mid           string
	Pid           string
	OverallThrows int
	Score         int
}
