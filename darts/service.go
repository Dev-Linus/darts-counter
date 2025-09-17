package darts

import (
	"errors"
	"log"

	playerstats "darts-counter/cmd/server/http/playerStats"
	playerthrow "darts-counter/cmd/server/http/playerThrow"
	models "darts-counter/models"
	response "darts-counter/response"
	storage "darts-counter/storage"
)

// Service is the service for the darts business logic.
type Service struct {
	Store    *storage.Storage
	Response response.Builder
}

// NewService creates a new darts Service.
func NewService(store *storage.Storage, resposneBuilder response.Builder) *Service {
	if store == nil || resposneBuilder == nil {
		log.Fatal("store or responseBuilder is nil")
	}
	return &Service{
		Store:    store,
		Response: resposneBuilder,
	}
}

// CollectStats aggregates statistics for the given player ID.
func (s *Service) CollectStats(_ string) (*playerstats.Response, error) {
	return nil, errors.New("error")
}

// PlayerThrow processes a player's throw in a match and returns the updated state.
func (s *Service) PlayerThrow(req *playerthrow.Request) (*playerthrow.Response, error) {
	if !isValidThrow(req.Throw) {
		return nil, errors.New("invalid throw")
	}

	mid := req.Mid
	pid := req.Pid

	match, err := s.Store.GetActiveMatch(mid)
	if err != nil {
		return nil, errors.New("error getting match or match is not active")
	}
	matchPlayerModel, err := s.Store.GetMatchPlayerModel(mid, pid)
	if err != nil {
		return nil, err
	}

	if matchPlayerModel.Score == match.StartAt { // is IN
		if !isValidIn(models.MapNumberToIO(match.StartMode), matchPlayerModel.Score, req.Throw) {
			// not a valid start, the turn is over, and it's the next players turn
			updatedMatch := s.persistTurnOver(match, req.Throw)

			return s.Response.BuildPlayerThrowResponse(updatedMatch, false, true), nil
		}

		updatedMatch, _, err := s.persistThrow(match, matchPlayerModel, &req.Throw)
		if err != nil {
			return nil, err
		}

		return s.Response.BuildPlayerThrowResponse(updatedMatch, false, false), nil
	}

	if matchPlayerModel.Score-req.Throw.ToPoints() == 0 { // is OUT
		if !isValidOut(models.MapNumberToIO(match.EndMode), matchPlayerModel.Score, req.Throw) {
			// build not valid out response
			updatedMatch := s.persistTurnOver(match, req.Throw)
			return s.Response.BuildPlayerThrowResponse(updatedMatch, false, true), nil
		}

		// valid finish, the player has won the game
		updatedMatch, _, err := s.persistThrow(match, matchPlayerModel, &req.Throw)
		if err != nil {
			return nil, err
		}
		// add a win to the player stats

		return s.Response.BuildPlayerThrowResponse(updatedMatch, true, false), nil
	}

	if isOverthrow(*match, matchPlayerModel.Score, req.Throw) {
		updatedMatch := s.persistTurnOver(match, req.Throw)

		return s.Response.BuildPlayerThrowResponse(updatedMatch, false, true), nil
	}
	// not IN not OUT not OVERTHROW => normal throw
	// persist normal throw
	// build persist response
	updatedMatch, _, err := s.persistThrow(match, matchPlayerModel, &req.Throw)
	if err != nil {
		return nil, err
	}

	return s.Response.BuildPlayerThrowResponse(updatedMatch, false, false), nil
}

func isValidThrow(throw models.ThrowType) bool {
	if _, ok := models.ThrowScores[throw]; !ok {
		return false
	}

	return true
}

func isValidIn(startMode models.IO, score int, throw models.ThrowType) bool {
	switch startMode {
	case models.Straight:
		return score-throw.ToPoints() > -1
	case models.Double:
		return throw.IsDouble() && score-throw.ToPoints() != 1
	case models.Master:
		return throw.IsMaster() && score-throw.ToPoints() != 1
	}
	return false
}

func isValidOut(endMode models.IO, _ int, throw models.ThrowType) bool {
	switch endMode {
	case models.Straight:
		return true
	case models.Double:
		return throw.IsDouble()
	case models.Master:
		return throw.IsMaster()
	}
	return false
}

func isOverthrow(match models.Match, score int, throw models.ThrowType) bool {
	endMode := models.MapNumberToIO(match.EndMode)
	potentialScore := score - throw.ToPoints()
	switch endMode {
	case models.Straight:
		return potentialScore < 0
	case models.Double:
		return potentialScore < 0 || potentialScore == 1
	case models.Master:
		return potentialScore < 0 || potentialScore == 1
	}
	return false
}

func (s *Service) persistTurnOver(match *models.Match, throw models.ThrowType) *models.Match {
	if _, err := s.Store.CreateThrow(match.ID, match.CurrentPlayer, int(throw)); err != nil {
		return nil
	}
	match.CurrentPlayer = match.GetNextPlayer()
	match.CurrentThrow = 0
	if err := s.Store.UpdateMatch(match); err != nil {
		return nil
	}

	return match
}

func (s *Service) persistThrow(match *models.Match, matchPlayerModel *models.MatchPlayer, throw *models.ThrowType) (*models.Match, *models.MatchPlayer, error) {
	match.Scores[match.CurrentPlayer] -= throw.ToPoints()
	match.CurrentThrow = (match.CurrentThrow + 1) % 3
	matchPlayerModel.OverallThrows++
	matchPlayerModel.Score = match.Scores[match.CurrentPlayer]
	if match.CurrentThrow == 0 && matchPlayerModel.Score > 0 {
		match.CurrentPlayer = match.GetNextPlayer()
	}

	_, err := s.Store.CreateThrow(matchPlayerModel.Mid, matchPlayerModel.Pid, int(*throw))
	if err != nil {
		return nil, nil, err
	}

	if matchPlayerModel.Score == 0 {
		if err := s.Store.WonMatch(match); err != nil {
			// Non-fatal: the response will still report a win; log for debugging
			log.Printf("warning: marking match %s as won failed: %v", match.ID, err)
		}
	}

	err = s.Store.UpdateMatch(match)
	if err != nil {
		return nil, nil, err
	}

	matchPlayerModel, err = s.Store.UpdateMatchPlayer(matchPlayerModel)
	if err != nil {
		return nil, nil, err
	}

	return match, matchPlayerModel, nil
}
