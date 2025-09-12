package models

// PlayerStats aggregates lifetime stats for a player across matches.
type PlayerStats struct {
	Pid        string `json:"pid"`
	Matches    int    `json:"matches"`
	Throws     int    `json:"throws"`
	TotalScore int    `json:"totalScore"`
}
