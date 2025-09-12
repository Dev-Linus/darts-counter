package storage

import (
	"database/sql"
	"errors"
	"log"

	"darts-counter/models"

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
		name TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS player_stats (
		pid TEXT PRIMARY KEY,
		matches INTEGER NOT NULL DEFAULT 0,
		throws INTEGER NOT NULL DEFAULT 0,
		totalScore INTEGER NOT NULL DEFAULT 0
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS matches (
		id TEXT PRIMARY KEY,
		isActive INTEGER NOT NULL DEFAULT 1,
		startAt INTEGER NOT NULL,
		startmode INTEGER NOT NULL,
		endmode INTEGER NOT NULL,
		currentPlayer TEXT,
		currentThrow INTEGER NOT NULL DEFAULT 0,
		wonBy TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS match_players (
		mid TEXT,
		pid TEXT,
		overallThrows INTEGER NOT NULL DEFAULT 0,
		score INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (mid, pid)
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS match_player_throws (
		id TEXT PRIMARY KEY,
		mid TEXT,
		pid TEXT,
		throwType INTEGER
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
		"INSERT INTO players (id, name) VALUES (?, ?)",
		id, name,
	)
	if err != nil {
		return nil, err
	}
	// initialize stats row
	_, _ = s.DB.Exec("INSERT OR IGNORE INTO player_stats (pid, matches, throws, totalScore) VALUES (?, 0, 0, 0)", id)

	return &models.Player{ID: id, Name: name}, nil
}

func (s *Storage) UpdatePlayer(id string, name string) (*models.Player, error) {
	_, err := s.DB.Exec("UPDATE players SET name=? WHERE id=?", name, id)
	if err != nil {
		return nil, err
	}
	return s.GetPlayer(id)
}

// UpdatePlayerModel updates by model and returns the updated instance
func (s *Storage) UpdatePlayerModel(p *models.Player) (*models.Player, error) {
	if p == nil || p.ID == "" {
		return nil, errors.New("invalid player model")
	}
	_, err := s.DB.Exec("UPDATE players SET name=? WHERE id=?", p.Name, p.ID)
	if err != nil {
		return nil, err
	}
	return s.GetPlayer(p.ID)
}

func (s *Storage) GetPlayers() ([]*models.Player, error) {
	rows, err := s.DB.Query("SELECT id, name FROM players")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var players []*models.Player
	for rows.Next() {
		var p models.Player
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		players = append(players, &p)
	}

	return players, nil
}

// GetAllPlayers returns all players (alias for GetPlayers for completeness)
func (s *Storage) GetAllPlayers() ([]*models.Player, error) {
	return s.GetPlayers()
}

func (s *Storage) GetPlayer(id string) (*models.Player, error) {
	row := s.DB.QueryRow("SELECT id, name FROM players WHERE id=?", id)
	var p models.Player
	if err := row.Scan(&p.ID, &p.Name); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Storage) DeletePlayer(id string) error {
	// remove throws for this player
	_, _ = s.DB.Exec("DELETE FROM match_player_throws WHERE pid=?", id)
	// remove match_player references for this player
	_, err := s.DB.Exec("DELETE FROM match_players WHERE pid=?", id)
	if err != nil {
		return err
	}
	// remove stats
	_, _ = s.DB.Exec("DELETE FROM player_stats WHERE pid=?", id)
	// remove player
	_, err = s.DB.Exec("DELETE FROM players WHERE id=?", id)
	return err
}

// ---------- MATCH METHODS ----------
func (s *Storage) CreateMatch(players []string, startAt int, startMode, endMode uint8) (*models.Match, error) {
	id := uuid.New().String()
	_, err := s.DB.Exec(
		"INSERT INTO matches (id, isActive, startAt, startmode, endmode, currentPlayer, currentThrow, wonBy) VALUES (?, 1, ?, ?, ?, '', 0, NULL)",
		id, startAt, startMode, endMode,
	)
	if err != nil {
		return nil, err
	}
	for _, pid := range players {
		_, err := s.DB.Exec("INSERT INTO match_players (mid, pid, overallThrows, score) VALUES (?, ?, 0, ?)", id, pid, startAt)
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
		Scores:       make(map[string]int),
	}
	for _, pid := range players {
		m.Scores[pid] = startAt
	}
	return m, nil
}

func (s *Storage) GetMatches() ([]*models.Match, error) {
	rows, err := s.DB.Query("SELECT id, startAt, startmode, endmode, currentThrow, currentPlayer FROM matches")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var matches []*models.Match
	for rows.Next() {
		var m models.Match
		if err := rows.Scan(&m.ID, &m.StartAt, &m.StartMode, &m.EndMode, &m.CurrentThrow, &m.CurrentPlayer); err != nil {
			return nil, err
		}

		m.Scores = make(map[string]int)
		pRows, err := s.DB.Query("SELECT pid, score FROM match_players WHERE mid=?", m.ID)
		if err != nil {
			return nil, err
		}
		for pRows.Next() {
			var pid string
			var score int
			if err = pRows.Scan(&pid, &score); err != nil {
				err = pRows.Close()
				if err != nil {
					return nil, err
				}
				return nil, err
			}
			m.Players = append(m.Players, pid)
			m.Scores[pid] = score
		}
		pRows.Close()

		matches = append(matches, &m)
	}
	return matches, nil
}

// GetAllMatches returns all matches (alias for GetMatches)
func (s *Storage) GetAllMatches() ([]*models.Match, error) {
	return s.GetMatches()
}

// GetMatch returns a single match by id
func (s *Storage) GetMatch(id string) (*models.Match, error) {
	rows, err := s.DB.Query("SELECT id, startAt, startmode, endmode, currentThrow, currentPlayer FROM matches WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	var m models.Match
	if err := rows.Scan(&m.ID, &m.StartAt, &m.StartMode, &m.EndMode, &m.CurrentThrow, &m.CurrentPlayer); err != nil {
		return nil, err
	}
	m.Scores = make(map[string]int)
	pRows, err := s.DB.Query("SELECT pid, score FROM match_players WHERE mid=?", m.ID)
	if err != nil {
		return nil, err
	}
	for pRows.Next() {
		var pid string
		var score int
		if err := pRows.Scan(&pid, &score); err != nil {
			pRows.Close()
			return nil, err
		}
		m.Players = append(m.Players, pid)
		m.Scores[pid] = score
	}
	pRows.Close()
	return &m, nil
}

func (s *Storage) UpdateMatch(match *models.Match) error {
	_, err := s.DB.Exec("UPDATE matches SET currentPlayer=?, currentThrow=? WHERE id=?", match.CurrentPlayer, match.CurrentThrow, match.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMatchModel updates a match and returns the updated instance
func (s *Storage) UpdateMatchModel(m *models.Match) (*models.Match, error) {
	if m == nil || m.ID == "" {
		return nil, errors.New("invalid match model")
	}
	_, err := s.DB.Exec("UPDATE matches SET currentPlayer=?, currentThrow=?, startAt=?, startmode=?, endmode=? WHERE id=?", m.CurrentPlayer, m.CurrentThrow, m.StartAt, m.StartMode, m.EndMode, m.ID)
	if err != nil {
		return nil, err
	}
	return s.GetMatch(m.ID)
}

func (s *Storage) DeleteMatch(id string) error {
	// remove throws for this match
	_, _ = s.DB.Exec("DELETE FROM match_player_throws WHERE mid=?", id)
	// remove match_players entries
	_, err := s.DB.Exec("DELETE FROM match_players WHERE mid=?", id)
	if err != nil {
		return err
	}
	// remove match
	_, err = s.DB.Exec("DELETE FROM matches WHERE id=?", id)
	return err
}

func (s *Storage) GetActiveMatch(mid string) (*models.Match, error) {
	rows, err := s.DB.Query("SELECT id, startAt, startmode, endmode, currentThrow, currentPlayer FROM matches WHERE id=? AND isActive=1", mid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	var m models.Match
	if err := rows.Scan(&m.ID, &m.StartAt, &m.StartMode, &m.EndMode, &m.CurrentThrow, &m.CurrentPlayer); err != nil {
		return nil, err
	}
	m.Scores = make(map[string]int)
	pRows, err := s.DB.Query("SELECT pid, score FROM match_players WHERE mid=?", m.ID)
	if err != nil {
		return nil, err
	}
	for pRows.Next() {
		var pid string
		var score int
		if err := pRows.Scan(&pid, &score); err != nil {
			pRows.Close()
			return nil, err
		}
		m.Players = append(m.Players, pid)
		m.Scores[pid] = score
	}
	pRows.Close()
	return &m, nil
}

// ---------- MATCH_PLAYER METHODS ----------
func (s *Storage) GetMatchPlayerModel(mid, pid string) (*models.MatchPlayer, error) {
	row := s.DB.QueryRow("SELECT mid, pid, overallThrows, score FROM match_players WHERE mid=? AND pid=?", mid, pid)
	if row == nil {
		return nil, errors.New("row is nil")
	}
	mp := &models.MatchPlayer{}
	if err := row.Scan(&mp.Mid, &mp.Pid, &mp.OverallThrows, &mp.Score); err != nil {
		return nil, err
	}
	return mp, nil
}

// ---------- GAMEPLAY ----------
func (s *Storage) RecordThrow(matchID, playerID string, amount int) (map[string]int, int, error) {
	// Update match player score and overall throws
	_, err := s.DB.Exec("UPDATE match_players SET score=score-?, overallThrows=overallThrows+1 WHERE mid=? AND pid=?", amount, matchID, playerID)
	if err != nil {
		return nil, 0, err
	}

	// Update player stats
	_, err = s.DB.Exec("UPDATE player_stats SET totalScore=totalScore+?, throws=throws+1 WHERE pid=?", amount, playerID)
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
	scoreRows, err := s.DB.Query("SELECT pid, score FROM match_players WHERE mid=?", matchID)
	if err != nil {
		return nil, 0, err
	}
	defer scoreRows.Close()
	scores := make(map[string]int)
	for scoreRows.Next() {
		var pid string
		var score int
		if err := scoreRows.Scan(&pid, &score); err != nil {
			return nil, 0, err
		}
		scores[pid] = score
	}
	return scores, nextThrow, nil
}
