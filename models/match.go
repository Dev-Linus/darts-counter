package models

type Match struct {
	ID            string         `json:"id"`
	Players       []string       `json:"players"`
	CurrentThrow  uint32         `json:"currentThrow"`
	CurrentPlayer string         `json:"currentPlayer"`
	StartAt       int            `json:"startAt"`
	StartMode     uint8          `json:"startMode"`
	EndMode       uint8          `json:"endMode"`
	Scores        map[string]int `json:"scores"`
}
