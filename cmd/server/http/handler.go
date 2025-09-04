package handler

import (
	creatematch "darts-counter/cmd/server/http/createMatch"
	createplayer "darts-counter/cmd/server/http/createPlayer"
	playerthrow "darts-counter/cmd/server/http/playerThrow"
	updateplayer "darts-counter/cmd/server/http/updatePlayer"
	"darts-counter/darts"
	"darts-counter/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Impl struct {
	Store        *storage.Storage
	DartsService *darts.Service
}

type Api interface {
	CreatePlayer(w http.ResponseWriter, r *http.Request)
	UpdatePlayer(w http.ResponseWriter, r *http.Request)
	DeletePlayer(w http.ResponseWriter, r *http.Request)
	ListPlayers(w http.ResponseWriter, r *http.Request)
	CreateMatch(w http.ResponseWriter, r *http.Request)
	ListMatches(w http.ResponseWriter, r *http.Request)
	DeleteMatch(w http.ResponseWriter, r *http.Request)
	PlayerThrow(w http.ResponseWriter, r *http.Request)
	Statistics(w http.ResponseWriter, r *http.Request)
	StreamFile(w http.ResponseWriter, r *http.Request)
}

func (i *Impl) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	req := &createplayer.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Name) < 1 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	p, err := i.Store.CreatePlayer(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func (i *Impl) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	req := &updateplayer.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(req.Name) < 1 {
		http.Error(w, "invalid name change requested", http.StatusBadRequest)
		return
	}

	if err := uuid.Validate(req.ID); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	p, err := i.Store.UpdatePlayer(req.ID, req.Name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func (i *Impl) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("playerId")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := i.Store.DeletePlayer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "player deleted"})
}

func (i *Impl) ListPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := i.Store.GetPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

func (i *Impl) CreateMatch(w http.ResponseWriter, r *http.Request) {
	req := &creatematch.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req == nil || len(req.Pids) < 1 {
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	for _, pid := range req.Pids {
		if err := uuid.Validate(pid); err != nil {
			http.Error(w, "invalid pid(s)", http.StatusBadRequest)
			return
		}
	}

	m, err := i.Store.CreateMatch(req.Pids, req.StartAt, req.StartMode, req.EndMode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(m)
}

func (i *Impl) ListMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := i.Store.GetMatches()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(matches)
}

func (s *Impl) DeleteMatch(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("matchId")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := s.Store.DeleteMatch(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "match deleted"})
}

func (i *Impl) PlayerThrow(w http.ResponseWriter, r *http.Request) {
	req := &playerthrow.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req == nil {
		http.Error(w, "req is nil", http.StatusBadRequest)
		return
	}

	if err := uuid.Validate(req.Pid); err != nil {
		http.Error(w, "invalid Pid", http.StatusBadRequest)
		return
	}

	if err := uuid.Validate(req.Mid); err != nil {
		http.Error(w, "invalid Mid", http.StatusBadRequest)
		return
	}

	resp, err := i.DartsService.PlayerThrow(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (i *Impl) Statistics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("playerId")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	playerStats, err := i.DartsService.CollectStats(id)
	if err != nil {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(playerStats)
}

func (i *Impl) StreamFile(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "missing file parameter", http.StatusBadRequest)
		return
	}

	// Restrict serving to your assets directory
	path := fmt.Sprintf("./assets/%s", file)

	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	// Detect MIME type (basic by extension)
	switch {
	case strings.HasSuffix(file, ".mp4"):
		w.Header().Set("Content-Type", "video/mp4")
	case strings.HasSuffix(file, ".mp3"):
		w.Header().Set("Content-Type", "audio/mpeg")
	default:
		// Fallback — browser will still try
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Stream with range support (important for video/audio seeking)
	http.ServeContent(w, r, file, time.Now(), f)
}

func NewApi(storage *storage.Storage, darts *darts.Service) (Api, error) {
	if storage == nil || darts == nil {
		return nil, errors.New("storage or darts service is nil")
	}

	return &Impl{
		Store:        storage,
		DartsService: darts,
	}, nil
}
