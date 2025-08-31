package models

type Player struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Matches    int    `json:"matches"`
	Throws     int    `json:"throws"`
	TotalScore int    `json:"totalScore"`
}
