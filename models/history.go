package models

type HistoryElement struct {
	Throw      ThrowType `json:"throw"`
	EndedTurn  bool      `json:"ended_turn"`
	TurnNumber int       `json:"turn_number"`
}

type History struct {
	History map[string][]HistoryElement `json:"history"`
}
