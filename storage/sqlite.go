package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"darts-counter/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

type Storage struct {
	Bun *bun.DB
}

// NewStorage creates a new SQLite service using Bun only
func NewStorage(dbFile string) *Storage {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	bunDB := bun.NewDB(db, sqlitedialect.New())
	ctx := context.Background()

	// Create tables with Bun (no migrations)
	if _, err := bunDB.NewCreateTable().Model((*playerRow)(nil)).IfNotExists().Exec(ctx); err != nil {
		log.Fatal(err)
	}
	if _, err := bunDB.NewCreateTable().Model((*playerStatsRow)(nil)).IfNotExists().Exec(ctx); err != nil {
		log.Fatal(err)
	}
	if _, err := bunDB.NewCreateTable().Model((*matchRow)(nil)).IfNotExists().Exec(ctx); err != nil {
		log.Fatal(err)
	}
	if _, err := bunDB.NewCreateTable().Model((*matchPlayerRow)(nil)).IfNotExists().Exec(ctx); err != nil {
		log.Fatal(err)
	}
	if _, err := bunDB.NewCreateTable().Model((*throwRow)(nil)).IfNotExists().Exec(ctx); err != nil {
		log.Fatal(err)
	}

	return &Storage{Bun: bunDB}
}

// ---------- PLAYER METHODS ----------

// CreatePlayer inserts a player using Bun ORM
func (s *Storage) CreatePlayer(name string) (*models.Player, error) {
	ctx := context.Background()
	p := &models.Player{ID: uuid.New().String(), Name: name}
	if _, err := s.Bun.NewInsert().Model(p).TableExpr("players").Exec(ctx); err != nil {
		return nil, err
	}
	// initialize stats row via Bun (insert ignore)
	ps := &playerStatsRow{Pid: p.ID}
	_, _ = s.Bun.NewInsert().Model(ps).On("CONFLICT (pid) DO NOTHING").Exec(ctx)
	return p, nil
}

// UpdatePlayer updates a player (Bun) and returns the updated record
func (s *Storage) UpdatePlayer(id string, name string) (*models.Player, error) {
	ctx := context.Background()
	// Update using Bun
	_, err := s.Bun.NewUpdate().TableExpr("players").Set("name = ?", name).Where("id = ?", id).Exec(ctx)
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
	ctx := context.Background()
	_, err := s.Bun.NewUpdate().TableExpr("players").Set("name = ?", p.Name).Where("id = ?", p.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetPlayer(p.ID)
}

func (s *Storage) GetPlayers() ([]*models.Player, error) {
	ctx := context.Background()
	var list []models.Player
	if err := s.Bun.NewSelect().TableExpr("players").Column("id", "name").Scan(ctx, &list); err != nil {
		return nil, err
	}
	players := make([]*models.Player, 0, len(list))
	for i := range list {
		p := list[i]
		players = append(players, &p)
	}
	return players, nil
}

// GetAllPlayers returns all players (alias for GetPlayers for completeness)
func (s *Storage) GetAllPlayers() ([]*models.Player, error) {
	return s.GetPlayers()
}

func (s *Storage) GetPlayer(id string) (*models.Player, error) {
	ctx := context.Background()
	var p models.Player
	if err := s.Bun.NewSelect().TableExpr("players").Column("id", "name").Where("id = ?", id).Scan(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Storage) DeletePlayer(id string) error {
	ctx := context.Background()
	_, _ = s.Bun.NewDelete().TableExpr("match_player_throws").Where("pid = ?", id).Exec(ctx)
	if _, err := s.Bun.NewDelete().TableExpr("match_players").Where("pid = ?", id).Exec(ctx); err != nil {
		return err
	}
	_, _ = s.Bun.NewDelete().TableExpr("player_stats").Where("pid = ?", id).Exec(ctx)
	_, err := s.Bun.NewDelete().TableExpr("players").Where("id = ?", id).Exec(ctx)
	return err
}

// ---------- MATCH METHODS ----------
func (s *Storage) CreateMatch(players []string, startAt int, startMode, endMode uint8) (*models.Match, error) {
	ctx := context.Background()
	id := uuid.New().String()
	mr := &matchRow{ID: id, IsActive: true, StartAt: startAt, Startmode: startMode, Endmode: endMode, CurrentPlayer: "", CurrentThrow: 0}
	if _, err := s.Bun.NewInsert().Model(mr).Exec(ctx); err != nil {
		return nil, err
	}
	for _, pid := range players {
		mpr := &matchPlayerRow{Mid: id, Pid: pid, OverallThrows: 0, Score: startAt}
		if _, err := s.Bun.NewInsert().Model(mpr).Exec(ctx); err != nil {
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
	ctx := context.Background()
	var mrows []matchRow
	if err := s.Bun.NewSelect().Model(&mrows).Column("id", "startAt", "startmode", "endmode", "currentThrow", "currentPlayer").Scan(ctx); err != nil {
		return nil, err
	}
	res := make([]*models.Match, 0, len(mrows))
	for _, mr := range mrows {
		m := &models.Match{
			ID:            mr.ID,
			Players:       []string{},
			CurrentThrow:  uint32(mr.CurrentThrow),
			CurrentPlayer: mr.CurrentPlayer,
			StartAt:       mr.StartAt,
			StartMode:     mr.Startmode,
			EndMode:       mr.Endmode,
			Scores:        make(map[string]int),
		}
		var mps []matchPlayerRow
		if err := s.Bun.NewSelect().Model(&mps).Column("pid", "score").Where("mid = ?", mr.ID).Scan(ctx); err != nil {
			return nil, err
		}
		for _, mp := range mps {
			m.Players = append(m.Players, mp.Pid)
			m.Scores[mp.Pid] = mp.Score
		}
		res = append(res, m)
	}
	return res, nil
}

// GetAllMatches returns all matches (alias for GetMatches)
func (s *Storage) GetAllMatches() ([]*models.Match, error) {
	return s.GetMatches()
}

// GetMatch returns a single match by id
func (s *Storage) GetMatch(id string) (*models.Match, error) {
	ctx := context.Background()
	var mr matchRow
	if err := s.Bun.NewSelect().Model(&mr).Column("id", "startAt", "startmode", "endmode", "currentThrow", "currentPlayer").Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	m := &models.Match{
		ID:            mr.ID,
		Players:       []string{},
		CurrentThrow:  uint32(mr.CurrentThrow),
		CurrentPlayer: mr.CurrentPlayer,
		StartAt:       mr.StartAt,
		StartMode:     mr.Startmode,
		EndMode:       mr.Endmode,
		Scores:        make(map[string]int),
	}
	var mps []matchPlayerRow
	if err := s.Bun.NewSelect().Model(&mps).Column("pid", "score").Where("mid = ?", mr.ID).Scan(ctx); err != nil {
		return nil, err
	}
	for _, mp := range mps {
		m.Players = append(m.Players, mp.Pid)
		m.Scores[mp.Pid] = mp.Score
	}
	return m, nil
}

func (s *Storage) UpdateMatch(match *models.Match) error {
	ctx := context.Background()
	_, err := s.Bun.NewUpdate().TableExpr("matches").
		Set("currentPlayer = ?", match.CurrentPlayer).
		Set("currentThrow = ?", match.CurrentThrow).
		Where("id = ?", match.ID).Exec(ctx)
	return err
}

// UpdateMatchModel updates a match and returns the updated instance
func (s *Storage) UpdateMatchModel(m *models.Match) (*models.Match, error) {
	if m == nil || m.ID == "" {
		return nil, errors.New("invalid match model")
	}
	ctx := context.Background()
	_, err := s.Bun.NewUpdate().TableExpr("matches").
		Set("currentPlayer = ?", m.CurrentPlayer).
		Set("currentThrow = ?", m.CurrentThrow).
		Set("startAt = ?", m.StartAt).
		Set("startmode = ?", m.StartMode).
		Set("endmode = ?", m.EndMode).
		Where("id = ?", m.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetMatch(m.ID)
}

func (s *Storage) DeleteMatch(id string) error {
	ctx := context.Background()
	// remove throws for this match
	_, _ = s.Bun.NewDelete().TableExpr("match_player_throws").Where("mid = ?", id).Exec(ctx)
	// remove match_players entries
	if _, err := s.Bun.NewDelete().TableExpr("match_players").Where("mid = ?", id).Exec(ctx); err != nil {
		return err
	}
	// remove match
	_, err := s.Bun.NewDelete().TableExpr("matches").Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Storage) GetActiveMatch(mid string) (*models.Match, error) {
	ctx := context.Background()
	var mr matchRow
	if err := s.Bun.NewSelect().Model(&mr).
		Column("id", "startAt", "startmode", "endmode", "currentThrow", "currentPlayer").
		Where("id = ?", mid).Where("isActive = 1").Scan(ctx); err != nil {
		return nil, err
	}
	m := &models.Match{
		ID:            mr.ID,
		Players:       []string{},
		CurrentThrow:  uint32(mr.CurrentThrow),
		CurrentPlayer: mr.CurrentPlayer,
		StartAt:       mr.StartAt,
		StartMode:     mr.Startmode,
		EndMode:       mr.Endmode,
		Scores:        make(map[string]int),
	}
	var mps []matchPlayerRow
	if err := s.Bun.NewSelect().Model(&mps).Column("pid", "score").Where("mid = ?", mr.ID).Scan(ctx); err != nil {
		return nil, err
	}
	for _, mp := range mps {
		m.Players = append(m.Players, mp.Pid)
		m.Scores[mp.Pid] = mp.Score
	}
	return m, nil
}

// ---------- MATCH_PLAYER METHODS ----------
func (s *Storage) GetMatchPlayerModel(mid, pid string) (*models.MatchPlayer, error) {
	ctx := context.Background()
	var mpr matchPlayerRow
	if err := s.Bun.NewSelect().Model(&mpr).Where("mid = ?", mid).Where("pid = ?", pid).Scan(ctx); err != nil {
		return nil, err
	}
	return &models.MatchPlayer{Mid: mpr.Mid, Pid: mpr.Pid, OverallThrows: mpr.OverallThrows, Score: mpr.Score}, nil
}

// ---------- GAMEPLAY ----------
func (s *Storage) RecordThrow(matchID, playerID string, amount int) (map[string]int, int, error) {
	ctx := context.Background()
	// Update match player score and overall throws
	if _, err := s.Bun.NewUpdate().TableExpr("match_players").
		Set("score = score - ?", amount).
		Set("overallThrows = overallThrows + 1").
		Where("mid = ?", matchID).Where("pid = ?", playerID).Exec(ctx); err != nil {
		return nil, 0, err
	}
	// Update player stats
	if _, err := s.Bun.NewUpdate().TableExpr("player_stats").
		Set("totalScore = totalScore + ?", amount).
		Set("throws = throws + 1").
		Where("pid = ?", playerID).Exec(ctx); err != nil {
		return nil, 0, err
	}
	// Update match current throw
	var currentThrow int
	if err := s.Bun.NewSelect().TableExpr("matches").Column("currentThrow").Where("id = ?", matchID).Scan(ctx, &currentThrow); err != nil {
		return nil, 0, err
	}
	nextThrow := currentThrow + 1
	if _, err := s.Bun.NewUpdate().TableExpr("matches").Set("currentThrow = ?", nextThrow).Where("id = ?", matchID).Exec(ctx); err != nil {
		return nil, 0, err
	}
	// Get updated scores
	var mps []matchPlayerRow
	if err := s.Bun.NewSelect().Model(&mps).Column("pid", "score").Where("mid = ?", matchID).Scan(ctx); err != nil {
		return nil, 0, err
	}
	scores := make(map[string]int)
	for _, mp := range mps {
		scores[mp.Pid] = mp.Score
	}
	return scores, nextThrow, nil
}

// ---- Bun table models (internal) ----
type playerRow struct {
	bun.BaseModel `bun:"table:players"`
	ID            string `bun:",pk"`
	Name          string `bun:",notnull"`
}

type playerStatsRow struct {
	bun.BaseModel `bun:"table:player_stats"`
	Pid           string `bun:",pk"`
	Matches       int    `bun:",notnull,default:0"`
	Throws        int    `bun:",notnull,default:0"`
	TotalScore    int    `bun:"totalScore,notnull,default:0"`
}

type matchRow struct {
	bun.BaseModel `bun:"table:matches"`
	ID            string  `bun:",pk"`
	IsActive      bool    `bun:"isActive,notnull,default:true"`
	StartAt       int     `bun:",notnull"`
	Startmode     uint8   `bun:"startmode,notnull"`
	Endmode       uint8   `bun:"endmode,notnull"`
	CurrentPlayer string  `bun:"currentPlayer,nullzero"`
	CurrentThrow  int     `bun:"currentThrow,notnull,default:0"`
	WonBy         *string `bun:"wonBy,nullzero"`
}

type matchPlayerRow struct {
	bun.BaseModel `bun:"table:match_players"`
	Mid           string `bun:",pk"`
	Pid           string `bun:",pk"`
	OverallThrows int    `bun:"overallThrows,notnull,default:0"`
	Score         int    `bun:",notnull,default:0"`
}

type throwRow struct {
	bun.BaseModel `bun:"table:match_player_throws"`
	ID            int64 `bun:",pk,autoincrement"`
	Mid           string
	Pid           string
	ThrowType     int
}
