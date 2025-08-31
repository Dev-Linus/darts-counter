package storage

import (
	"darts-counter/models"
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(dbFile string) *Storage {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS players (
		id TEXT PRIMARY KEY,
		name TEXT,
		matches INTEGER,
		throws INTEGER,
		totalScore INTEGER
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS matches (
		id TEXT PRIMARY KEY,
		startAt INTEGER,
		startmode CHAR,
		endmode CHAR,
		currentThrow INTEGER,
		currentScore INTEGER
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS match_players (
		matchId TEXT,
		playerId TEXT,
		throws INTEGER,
		score INTEGER,
		PRIMARY KEY (matchId, playerId)
	);`)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{DB: db}
}

// ---------- PLAYER METHODS ----------
func (s *Storage) CreatePlayer(name string) (*models.Player, error) {
	id := uuid.New().String()
	_, err := s.DB.Exec(
		"INSERT INTO players (id, name, matches, throws, totalScore) VALUES (?, ?, ?, ?, ?)",
		id, name, 0, 0, 0,
	)
	if err != nil {
		return nil, err
	}

	return &models.Player{ID: id, Name: name, Matches: 0, Throws: 0, TotalScore: 0}, nil
}

func (s *Storage) UpdatePlayer(id string, name string) (*models.Player, error) {
	_, err := s.DB.Exec("UPDATE players SET name=? WHERE id=?", name, id)
	if err != nil {
		return nil, err
	}

	return &models.Player{ID: id, Name: name}, nil
}

func (s *Storage) GetPlayers() ([]*models.Player, error) {
	rows, err := s.DB.Query("SELECT id, name, matches, throws, totalScore FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*models.Player
	for rows.Next() {
		var p models.Player
		rows.Scan(&p.ID, &p.Name, &p.Matches, &p.Throws, &p.TotalScore)
		players = append(players, &p)
	}

	return players, nil
}

func (s *Storage) GetPlayer(id string) (*models.Player, error) {
	row := s.DB.QueryRow("SELECT id, name, matches, throws, totalScore FROM players WHERE id=?", id)
	var p models.Player
	err := row.Scan(&p.ID, &p.Name, &p.Matches, &p.Throws, &p.TotalScore)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *Storage) DeletePlayer(id string) error {
	_, err := s.DB.Exec("DELETE FROM match_players WHERE playerId=?", id)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec("DELETE FROM players WHERE id=?", id)

	return err
}

// ---------- MATCH METHODS ----------
func (s *Storage) CreateMatch(players []string, startAt int, startMode, endMode uint8) (*models.Match, error) {
	id := uuid.New().String()
	_, err := s.DB.Exec(
		"INSERT INTO matches (id, startAt, startmode, endmode, currentThrow, currentScore) VALUES (?, ?, ?, ?, ?, ?)",
		id, startAt, startMode, endMode, 0, startAt,
	)
	if err != nil {
		return nil, err
	}
	for _, pid := range players {
		_, err := s.DB.Exec("INSERT INTO match_players (matchId, playerId, throws, score) VALUES (?, ?, ?, ?)", id, pid, 0, 0)
		if err != nil {
			return nil, err
		}
	}

	m := &models.Match{
		ID:           id,
		Players:      players,
		CurrentThrow: 0,
		StartAt:      startAt,
		StartMode:    startMode,
		EndMode:      endMode,
		CurrentScore: startAt,
		Scores:       make(map[string]int),
	}
	for _, pid := range players {
		m.Scores[pid] = 0
	}
	return m, nil
}

func (s *Storage) GetMatches() ([]*models.Match, error) {
	rows, err := s.DB.Query("SELECT id, startAt, startMode, endMode, currentThrow, currentScore FROM matches")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*models.Match
	for rows.Next() {
		var m models.Match
		rows.Scan(&m.ID, &m.StartAt, &m.StartMode, &m.EndMode, &m.CurrentThrow, &m.CurrentScore)

		m.Scores = make(map[string]int)
		pRows, _ := s.DB.Query("SELECT playerId, score FROM match_players WHERE matchId=?", m.ID)
		for pRows.Next() {
			var pid string
			var score int
			pRows.Scan(&pid, &score)
			m.Players = append(m.Players, pid)
			m.Scores[pid] = score
		}
		pRows.Close()

		matches = append(matches, &m)
	}
	return matches, nil
}

func (s *Storage) DeleteMatch(id string) error {
	_, err := s.DB.Exec("DELETE FROM match_players WHERE matchId=?", id)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec("DELETE FROM matches WHERE id=?", id)
	return err
}

// ---------- GAMEPLAY ----------
func (s *Storage) RecordThrow(matchID, playerID string, amount int) (map[string]int, int, error) {
	// Update player stats
	_, err := s.DB.Exec("UPDATE players SET totalScore=totalScore+?, throws=throws+1 WHERE id=?", amount, playerID)
	if err != nil {
		return nil, 0, err
	}

	// Update match player score
	_, err = s.DB.Exec("UPDATE match_players SET score=score+?, throws=throws+1 WHERE matchId=? AND playerId=?", amount, matchID, playerID)
	if err != nil {
		return nil, 0, err
	}

	// Update match current throw
	row := s.DB.QueryRow("SELECT currentThrow FROM matches WHERE id=?", matchID)
	var currentThrow int
	if err := row.Scan(&currentThrow); err != nil {
		return nil, 0, err
	}
	nextThrow := currentThrow + 1
	_, err = s.DB.Exec("UPDATE matches SET currentThrow=? WHERE id=?", nextThrow, matchID)
	if err != nil {
		return nil, 0, err
	}

	// Get updated scores
	scoreRows, _ := s.DB.Query("SELECT playerId, score FROM match_players WHERE matchId=?", matchID)
	defer scoreRows.Close()
	scores := make(map[string]int)
	for scoreRows.Next() {
		var pid string
		var score int
		scoreRows.Scan(&pid, &score)
		scores[pid] = score
	}
	return scores, nextThrow, nil
}
