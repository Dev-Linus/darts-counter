package storage

import (
	"context"
	"errors"

	"darts-counter/models"
)

// Additional CRUD coverage for remaining tables and convenience helpers

// CreatePlayerStatsDefault
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

// CreateMatchPlayer
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

// ThrowRecord is a minimal model for match_player_throws table
type ThrowRecord struct {
	ID        int64
	Mid       string
	Pid       string
	ThrowType int // use models.ThrowType values
	EndedTurn bool
	Turn      int
}

func (s *Storage) CreateThrow(tr ThrowRecord) (*ThrowRecord, error) {
	if tr.Mid == "" || tr.Pid == "" {
		return nil, errors.New("empty ids")
	}
	ctx := context.Background()
	count, err := s.countEndedTurns(ctx, tr.Mid, tr.Pid)
	if err != nil {
		return nil, err
	}
	tr.Turn = 1 + count

	row := &throwRow{Mid: tr.Mid, Pid: tr.Pid, ThrowType: tr.ThrowType, EndedTurn: tr.EndedTurn, Turn: tr.Turn}
	if err := s.Bun.NewInsert().Model(row).Returning("*").Scan(ctx); err != nil {
		return nil, err
	}
	return &ThrowRecord{ID: row.ID, Mid: row.Mid, Pid: row.Pid, ThrowType: row.ThrowType, EndedTurn: row.EndedTurn, Turn: row.Turn}, nil
}

func (s *Storage) countEndedTurns(ctx context.Context, mid, pid string) (int, error) {
	return s.Bun.
		NewSelect().
		Model((*throwRow)(nil)). // no allocation; just use the table
		Where("mid = ?", mid).
		Where("pid = ?", pid).
		Where("endedTurn = ?", true). // or .Where("endedTurn = 1") for SQLite
		Count(ctx)
}

func (s *Storage) GetThrow(id int64) (*ThrowRecord, error) {
	ctx := context.Background()
	var r throwRow
	if err := s.Bun.NewSelect().Model(&r).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &ThrowRecord{ID: r.ID, Mid: r.Mid, Pid: r.Pid, ThrowType: r.ThrowType, EndedTurn: r.EndedTurn, Turn: r.Turn}, nil
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
		Set("endedTurn = ?", tr.EndedTurn).
		Set("turn = ?", tr.Turn).
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
