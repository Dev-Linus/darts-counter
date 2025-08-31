package darts

import (
	"errors"

	playerstats "darts-counter/cmd/server/http/playerStats"
	playerthrow "darts-counter/cmd/server/http/playerThrow"
	"darts-counter/storage"
)

type Service struct {
	Store *storage.Storage
}

func NewService(store *storage.Storage) *Service {
	return &Service{Store: store}
}

func (s *Service) CollectStats(pid string) (*playerstats.Response, error) {
	return nil, errors.New("error")
}

func (s *Service) PlayerThrow(req *playerthrow.Request) (*playerthrow.Response, error) {
	return nil, errors.New("error")
}
