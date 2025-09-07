package darts

import (
	"errors"

	playerstats "darts-counter/cmd/server/http/playerStats"
	playerthrow "darts-counter/cmd/server/http/playerThrow"
	"darts-counter/models"
	"darts-counter/response"
	"darts-counter/storage"
)

type Service struct {
	Store    *storage.Storage
	Response response.Builder
}

func NewService(store *storage.Storage, resposneBuilder response.Builder) *Service {
	return &Service{
		Store:    store,
		Response: resposneBuilder,
	}
}

func (s *Service) CollectStats(pid string) (*playerstats.Response, error) {
	return nil, errors.New("error")
}

func (s *Service) PlayerThrow(req *playerthrow.Request) (*playerthrow.Response, error) {
	if !isValidThrow(req.Throw) {
		return nil, errors.New("invalid throw")
	}

	mid := req.Mid
	pid := req.Pid

	match, err := s.Store.GetActiveMatch(mid)
	if err != nil {
		return nil, err
	}
	matchPlayerModel, err := s.Store.GetMatchPlayerModel(mid, pid)
	if err != nil {
		return nil, err
	}

	if matchPlayerModel.Score == match.StartAt { // is IN
		if !isValidIn(models.MapNumberToIO(match.StartMode), matchPlayerModel.Score, req.Throw) {
			// not a valid start, the turn over and its the next players turn
			updatedMatch := s.persistTurnOver(match)

			return s.Response.BuildNextPlayerResponse(updatedMatch), errors.New("no valid start throw")
		}

		// persistThrow return error
		// build persist response

		/*if err := persistThrow(match, matchPlayerModel, req); err != nil {
			// db error
		}

		resp, err := s.Response.BuildPersistPlayerThrowResponse(match, matchPlayerModel, req)
		if err != nil {
			// error
		}*/

		/*scores := s.Store.PersistThrow(mid, pid, req.Throw)
		resp := &playerthrow.Response{
			Won:    false,
			Scores: scores,
		}

		if match.CurrentThrow >= 2 {
			resp.NextThrowBy = match.NextPlayer()
			resp.PossibleFinish = computePossibleFinishForPlayer(match.NextPlayer())
			s.Store.UpdateThrowsThisTurn(mid, pid, 0)
		}

		resp.NextThrowBy = pid
		s.Store.UpdateThrowsThisTurn(mid, pid, matchPlayerModel.throwsThisTurn+1)
		resp.PossibleFinish = computePossibleFinishForPlayer(pid)*/
	}

	/*if matchPlayerModel.Score-req.Throw.ToPoints() == 0 { // is OUT
		if !isValidOut(models.MapNumberToIO(match.EndMode), matchPlayerModel.Score, req.Throw) {
			// build not valid out response
			return &playerthrow.Response{
				Won: false,
			}, errors.New("no valid end throw")
		}

		// valid finish, player has won the game
		//if err := persistThrow(match, matchPlayerModel, req); err != nil {
			// db error
		//}
		// add a win to the player stats
		// deactivate the match
		// compute scores

		return s.Response.BuildPersistPlayerThrowResponse(match, matchPlayerModel, req)
	}*/

	if isOverthrow(*match, matchPlayerModel.Score, req.Throw) { // not IN not OUT but OVERTHROW
		// build overthrow response
		/*scores := s.Store.GetScoresForMatch(mid)
		return &playerthrow.Response{
			Won:            false,
			NextThrowBy:    match.NextPlayer(),
			Scores:         scores,
			PossibleFinish: computePossibleFinishForPlayer(match.NextPlayer()),
		}, errors.New("overthrown, next players turn")*/
	}

	// not IN not OUT not OVERTHROW => normal throw
	// persist normal throw
	// build persist resposne

	/*if err := persistThrow(match, matchPlayerModel, req); err != nil {
		// db error
	}*/

	/*resp, err := s.Response.BuildPersistPlayerThrowResponse(match, matchPlayerModel, req)
	if err != nil {
		// error
	}*/

	return nil, nil
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

func isValidOut(endMode models.IO, score int, throw models.ThrowType) bool {
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
	throwsLeft := 3 - match.CurrentThrow
	endMode := models.MapNumberToIO(match.EndMode)
	potentialScore := score - throw.ToPoints()
	switch endMode {
	case models.Straight:
		return potentialScore < 0
	case models.Double:
		return potentialScore == 1 || (throwsLeft == 1 && potentialScore%2 != 0)
	case models.Master:
		return potentialScore == 1 || (throwsLeft == 1 && potentialScore%3 != 0 && potentialScore%2 != 0)
	}
	return false
}

func (s Service) persistTurnOver(match *models.Match) *models.Match {
	match.CurrentPlayer = match.GetNextPlayer()
	match.CurrentThrow = 1
	if err := s.Store.UpdateMatch(match); err != nil {
		return nil
	}

	return match
}

/*func persistThrow() {
	scores := s.Store.PersistThrow(mid, pid, req.Throw)

	resp := &playerthrow.Response{
		Won:    false,
		Scores: scores,
	}

	if matchPlayerModel.throwsThisTurn >= 2 {
		resp.NextThrowBy = match.NextPlayer()
		resp.PossibleFinish = computePossibleFinishForPlayer(match.NextPlayer())
		s.Store.UpdateThrowsThisTurn(mid, pid, 0)
	}

	resp.NextThrowBy = pid
	s.Store.UpdateThrowsThisTurn(mid, pid, matchPlayerModel.throwsThisTurn+1)
	resp.PossibleFinish = computePossibleFinishForPlayer(pid)

	return resp, nil
}*/
