package models

type PlayerStats struct {
	Pid        string `json:"pid"`
	Matches    int    `json:"matches"`
	Throws     int    `json:"throws"`
	TotalScore int    `json:"totalScore"`
}
