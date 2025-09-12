package storage

import (
	"errors"

	"github.com/google/uuid"

	"darts-counter/models"
)

// Additional CRUD coverage for remaining tables and convenience helpers

// ---------- PLAYER_STATS CRUD ----------
func (s *Storage) CreatePlayerStatsDefault(pid string) (*models.PlayerStats, error) {
	if pid == "" {
		return nil, errors.New("empty pid")
	}
	_, err := s.DB.Exec("INSERT OR IGNORE INTO player_stats (pid, matches, throws, totalScore) VALUES (?, 0, 0, 0)", pid)
	if err != nil {
		return nil, err
	}
	return s.GetPlayerStats(pid)
}

func (s *Storage) GetPlayerStats(pid string) (*models.PlayerStats, error) {
	row := s.DB.QueryRow("SELECT pid, matches, throws, totalScore FROM player_stats WHERE pid=?", pid)
	ps := &models.PlayerStats{}
	if err := row.Scan(&ps.Pid, &ps.Matches, &ps.Throws, &ps.TotalScore); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *Storage) GetAllPlayerStats() ([]*models.PlayerStats, error) {
	rows, err := s.DB.Query("SELECT pid, matches, throws, totalScore FROM player_stats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.PlayerStats
	for rows.Next() {
		ps := &models.PlayerStats{}
		if err := rows.Scan(&ps.Pid, &ps.Matches, &ps.Throws, &ps.TotalScore); err != nil {
			return nil, err
		}
		out = append(out, ps)
	}
	return out, nil
}

func (s *Storage) UpdatePlayerStats(ps *models.PlayerStats) (*models.PlayerStats, error) {
	if ps == nil || ps.Pid == "" {
		return nil, errors.New("invalid player stats model")
	}
	_, err := s.DB.Exec("UPDATE player_stats SET matches=?, throws=?, totalScore=? WHERE pid=?", ps.Matches, ps.Throws, ps.TotalScore, ps.Pid)
	if err != nil {
		return nil, err
	}
	return s.GetPlayerStats(ps.Pid)
}

func (s *Storage) DeletePlayerStats(pid string) error {
	_, err := s.DB.Exec("DELETE FROM player_stats WHERE pid=?", pid)
	return err
}

// ---------- MATCH_PLAYER CRUD ----------
func (s *Storage) CreateMatchPlayer(mid, pid string, startAt int) (*models.MatchPlayer, error) {
	if mid == "" || pid == "" {
		return nil, errors.New("empty ids")
	}
	_, err := s.DB.Exec("INSERT INTO match_players (mid, pid, overallThrows, score) VALUES (?, ?, 0, ?)", mid, pid, startAt)
	if err != nil {
		return nil, err
	}
	return s.GetMatchPlayerModel(mid, pid)
}

func (s *Storage) GetAllMatchPlayers(mid string) ([]*models.MatchPlayer, error) {
	rows, err := s.DB.Query("SELECT mid, pid, overallThrows, score FROM match_players WHERE mid=?", mid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.MatchPlayer
	for rows.Next() {
		mp := &models.MatchPlayer{}
		if err := rows.Scan(&mp.Mid, &mp.Pid, &mp.OverallThrows, &mp.Score); err != nil {
			return nil, err
		}
		out = append(out, mp)
	}
	return out, nil
}

func (s *Storage) UpdateMatchPlayer(mp *models.MatchPlayer) (*models.MatchPlayer, error) {
	if mp == nil || mp.Mid == "" || mp.Pid == "" {
		return nil, errors.New("invalid match player model")
	}
	_, err := s.DB.Exec("UPDATE match_players SET overallThrows=?, score=? WHERE mid=? AND pid=?", mp.OverallThrows, mp.Score, mp.Mid, mp.Pid)
	if err != nil {
		return nil, err
	}
	return s.GetMatchPlayerModel(mp.Mid, mp.Pid)
}

func (s *Storage) DeleteMatchPlayer(mid, pid string) error {
	_, err := s.DB.Exec("DELETE FROM match_players WHERE mid=? AND pid=?", mid, pid)
	return err
}

// ---------- MATCH_PLAYER_THROWS CRUD ----------
// ThrowRecord is a minimal model for match_player_throws table
// We keep it local to storage to avoid changing public models if not needed.
type ThrowRecord struct {
	ID        string
	Mid       string
	Pid       string
	ThrowType int // use models.ThrowType values
}

func (s *Storage) CreateThrow(mid, pid string, throwType int) (*ThrowRecord, error) {
	if mid == "" || pid == "" {
		return nil, errors.New("empty ids")
	}
	tid := uuid.New().String()
	_, err := s.DB.Exec("INSERT INTO match_player_throws (id, mid, pid, throwType) VALUES (?, ?, ?, ?)", tid, mid, pid, throwType)
	if err != nil {
		return nil, err
	}
	return &ThrowRecord{ID: tid, Mid: mid, Pid: pid, ThrowType: throwType}, nil
}

func (s *Storage) GetThrow(id string) (*ThrowRecord, error) {
	row := s.DB.QueryRow("SELECT id, mid, pid, throwType FROM match_player_throws WHERE id=?", id)
	tr := &ThrowRecord{}
	if err := row.Scan(&tr.ID, &tr.Mid, &tr.Pid, &tr.ThrowType); err != nil {
		return nil, err
	}
	return tr, nil
}

func (s *Storage) GetAllThrowsForMatchPlayer(mid, pid string) ([]*ThrowRecord, error) {
	rows, err := s.DB.Query("SELECT id, mid, pid, throwType FROM match_player_throws WHERE mid=? AND pid=?", mid, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ThrowRecord
	for rows.Next() {
		tr := &ThrowRecord{}
		if err := rows.Scan(&tr.ID, &tr.Mid, &tr.Pid, &tr.ThrowType); err != nil {
			return nil, err
		}
		out = append(out, tr)
	}
	return out, nil
}

func (s *Storage) UpdateThrow(tr *ThrowRecord) (*ThrowRecord, error) {
	if tr == nil || tr.ID == "" {
		return nil, errors.New("invalid throw model")
	}
	_, err := s.DB.Exec("UPDATE match_player_throws SET mid=?, pid=?, throwType=? WHERE id=?", tr.Mid, tr.Pid, tr.ThrowType, tr.ID)
	if err != nil {
		return nil, err
	}
	return s.GetThrow(tr.ID)
}

func (s *Storage) DeleteThrow(id string) error {
	_, err := s.DB.Exec("DELETE FROM match_player_throws WHERE id=?", id)
	return err
}
