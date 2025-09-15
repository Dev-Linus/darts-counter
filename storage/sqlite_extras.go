package storage

import (
	"context"
	"errors"

	"darts-counter/models"
)

// Additional CRUD coverage for remaining tables and convenience helpers

// ---------- PLAYER_STATS CRUD ----------
func (s *Storage) CreatePlayerStatsDefault(pid string) (*models.PlayerStats, error) {
	if pid == "" {
		return nil, errors.New("empty pid")
	}
	ctx := context.Background()
	ps := &playerStatsRow{Pid: pid}
	_, err := s.Bun.NewInsert().Model(ps).On("CONFLICT (pid) DO NOTHING").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetPlayerStats(pid)
}

func (s *Storage) GetPlayerStats(pid string) (*models.PlayerStats, error) {
	ctx := context.Background()
	var pr playerStatsRow
	if err := s.Bun.NewSelect().Model(&pr).Where("pid = ?", pid).Scan(ctx); err != nil {
		return nil, err
	}
	return &models.PlayerStats{Pid: pr.Pid, Matches: pr.Matches, Throws: pr.Throws, TotalScore: pr.TotalScore}, nil
}

// GetAllPlayerStats
func (s *Storage) GetAllPlayerStats() ([]*models.PlayerStats, error) {
	ctx := context.Background()
	var rows []playerStatsRow
	if err := s.Bun.NewSelect().Model(&rows).Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]*models.PlayerStats, 0, len(rows))
	for _, r := range rows {
		out = append(out, &models.PlayerStats{Pid: r.Pid, Matches: r.Matches, Throws: r.Throws, TotalScore: r.TotalScore})
	}
	return out, nil
}

func (s *Storage) UpdatePlayerStats(ps *models.PlayerStats) (*models.PlayerStats, error) {
	if ps == nil || ps.Pid == "" {
		return nil, errors.New("invalid player stats model")
	}
	ctx := context.Background()
	_, err := s.Bun.NewUpdate().TableExpr("player_stats").
		Set("matches = ?", ps.Matches).
		Set("throws = ?", ps.Throws).
		Set("totalScore = ?", ps.TotalScore).
		Where("pid = ?", ps.Pid).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetPlayerStats(ps.Pid)
}

func (s *Storage) DeletePlayerStats(pid string) error {
	ctx := context.Background()
	_, err := s.Bun.NewDelete().TableExpr("player_stats").Where("pid = ?", pid).Exec(ctx)
	return err
}

// ---------- MATCH_PLAYER CRUD ----------
func (s *Storage) CreateMatchPlayer(mid, pid string, startAt int) (*models.MatchPlayer, error) {
	if mid == "" || pid == "" {
		return nil, errors.New("empty ids")
	}
	ctx := context.Background()
	mpr := &matchPlayerRow{Mid: mid, Pid: pid, OverallThrows: 0, Score: startAt}
	if _, err := s.Bun.NewInsert().Model(mpr).Exec(ctx); err != nil {
		return nil, err
	}
	return s.GetMatchPlayerModel(mid, pid)
}

func (s *Storage) GetAllMatchPlayers(mid string) ([]*models.MatchPlayer, error) {
	ctx := context.Background()
	var rows []matchPlayerRow
	if err := s.Bun.NewSelect().Model(&rows).Where("mid = ?", mid).Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]*models.MatchPlayer, 0, len(rows))
	for _, r := range rows {
		out = append(out, &models.MatchPlayer{Mid: r.Mid, Pid: r.Pid, OverallThrows: r.OverallThrows, Score: r.Score})
	}
	return out, nil
}

func (s *Storage) UpdateMatchPlayer(mp *models.MatchPlayer) (*models.MatchPlayer, error) {
	if mp == nil || mp.Mid == "" || mp.Pid == "" {
		return nil, errors.New("invalid match player model")
	}
	ctx := context.Background()
	_, err := s.Bun.NewUpdate().TableExpr("match_players").
		Set("overallThrows = ?", mp.OverallThrows).
		Set("score = ?", mp.Score).
		Where("mid = ?", mp.Mid).Where("pid = ?", mp.Pid).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetMatchPlayerModel(mp.Mid, mp.Pid)
}

func (s *Storage) DeleteMatchPlayer(mid, pid string) error {
	ctx := context.Background()
	_, err := s.Bun.NewDelete().TableExpr("match_players").Where("mid = ?", mid).Where("pid = ?", pid).Exec(ctx)
	return err
}

// ---------- MATCH_PLAYER_THROWS CRUD ----------
// ThrowRecord is a minimal model for match_player_throws table
// We keep it local to storage to avoid changing public models if not needed.
type ThrowRecord struct {
	ID        int64
	Mid       string
	Pid       string
	ThrowType int // use models.ThrowType values
}

func (s *Storage) CreateThrow(mid, pid string, throwType int) (*ThrowRecord, error) {
	if mid == "" || pid == "" {
		return nil, errors.New("empty ids")
	}
	ctx := context.Background()
	row := &throwRow{Mid: mid, Pid: pid, ThrowType: throwType}
	res, err := s.Bun.NewInsert().Model(row).Exec(ctx)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &ThrowRecord{ID: id, Mid: mid, Pid: pid, ThrowType: throwType}, nil
}

func (s *Storage) GetThrow(id int64) (*ThrowRecord, error) {
	ctx := context.Background()
	var r throwRow
	if err := s.Bun.NewSelect().Model(&r).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &ThrowRecord{ID: r.ID, Mid: r.Mid, Pid: r.Pid, ThrowType: r.ThrowType}, nil
}

func (s *Storage) GetAllThrowsForMatchPlayer(mid, pid string) ([]*ThrowRecord, error) {
	ctx := context.Background()
	var rows []throwRow
	if err := s.Bun.NewSelect().Model(&rows).Where("mid = ?", mid).Where("pid = ?", pid).Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]*ThrowRecord, 0, len(rows))
	for _, r := range rows {
		out = append(out, &ThrowRecord{ID: r.ID, Mid: r.Mid, Pid: r.Pid, ThrowType: r.ThrowType})
	}
	return out, nil
}

func (s *Storage) UpdateThrow(tr *ThrowRecord) (*ThrowRecord, error) {
	if tr == nil || tr.ID == 0 {
		return nil, errors.New("invalid throw model")
	}
	ctx := context.Background()
	_, err := s.Bun.NewUpdate().TableExpr("match_player_throws").
		Set("mid = ?", tr.Mid).
		Set("pid = ?", tr.Pid).
		Set("throwType = ?", tr.ThrowType).
		Where("id = ?", tr.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetThrow(tr.ID)
}

func (s *Storage) DeleteThrow(id int64) error {
	ctx := context.Background()
	_, err := s.Bun.NewDelete().TableExpr("match_player_throws").Where("id = ?", id).Exec(ctx)
	return err
}

// -------- High-level helpers that accept Match and ThrowType models (CRD) --------
// CreateMatchPlayerThrow creates a throw for the given match model, player id and typed throw.
func (s *Storage) CreateMatchPlayerThrow(m *models.Match, pid string, tt models.ThrowType) (*ThrowRecord, error) {
	if m == nil || m.ID == "" {
		return nil, errors.New("invalid match model")
	}
	return s.CreateThrow(m.ID, pid, int(tt))
}

// GetMatchPlayerThrows returns all throws for the given match model and player id.
func (s *Storage) GetMatchPlayerThrows(m *models.Match, pid string) ([]*ThrowRecord, error) {
	if m == nil || m.ID == "" {
		return nil, errors.New("invalid match model")
	}
	return s.GetAllThrowsForMatchPlayer(m.ID, pid)
}

// DeleteMatchPlayerThrow deletes a specific throw (by id) for the given match model.
// The match is accepted to satisfy the API requirement; its ID is not strictly needed for deletion but helps validate context.
func (s *Storage) DeleteMatchPlayerThrow(m *models.Match, throwID int64) error {
	if m == nil || m.ID == "" {
		return errors.New("invalid match model")
	}
	if throwID == 0 {
		return errors.New("empty throw id")
	}
	return s.DeleteThrow(throwID)
}
