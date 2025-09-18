package models

// Match represents a darts match state.
type Match struct {
	ID            string         `json:"id"`
	Players       []string       `json:"players"`
	CurrentThrow  uint32         `json:"currentThrow"`
	CurrentPlayer string         `json:"currentPlayer"`
	WonBy         string         `json:"wonBy"`
	StartAt       int            `json:"startAt"`
	StartMode     uint8          `json:"startMode"`
	EndMode       uint8          `json:"endMode"`
	Scores        map[string]int `json:"scores"`
}

// GetNextPlayer returns the next player's ID in the rotation.
func (m *Match) GetNextPlayer() string {
	for i, pid := range m.Players {
		if pid != m.CurrentPlayer {
			continue
		}

		return m.Players[(i+1)%len(m.Players)]
	}

	return ""
}
