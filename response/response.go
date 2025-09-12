package response

import (
	playerthrow "darts-counter/cmd/server/http/playerThrow"
	"darts-counter/models"
)

type Builder interface {
	BuildPlayerResponse(match *models.Match, won bool) *playerthrow.Response
}

type Impl struct {
}

func NewBuilder() Builder {
	return Impl{}
}

func (i Impl) BuildPlayerResponse(match *models.Match, won bool) *playerthrow.Response {
	return &playerthrow.Response{
		Won:            won,
		NextThrowBy:    match.CurrentPlayer,
		Scores:         match.Scores,
		PossibleFinish: getPossibleFinishForMatchPlayer(match),
	}
}

func (i Impl) BuildPersistPlayerThrowResponse(match *models.Match, currentPid string) *playerthrow.Response {
	return nil
}

func getPossibleFinishForMatchPlayer(match *models.Match) []models.ThrowType {
	playerScore, ok := match.Scores[match.CurrentPlayer]
	if !ok {
		return nil
	}
	throwsLeft := 3 - int(match.CurrentThrow)
	endMode := models.MapNumberToIO(match.EndMode)
	throws := float32(playerScore) / float32(60)
	if float32(throwsLeft) < throws {
		return nil
	}

	bestFinish := make([]models.ThrowType, 0, throwsLeft)
	possibleLastThrows := endMode.GetAllFinishingThrows()

	for _, possibleThrow := range possibleLastThrows {
		potentialScore := playerScore - possibleThrow.ToPoints()
		if potentialScore == 0 {
			bestFinish = append(bestFinish, possibleThrow)
			break
		}
		if potentialScore > 0 && throwsLeft-1 > 0 {
			followingThrows := make([]models.ThrowType, 0, throwsLeft)
			followingThrows = computePossibleNextThrows(potentialScore, throwsLeft-1, followingThrows)
			if len(followingThrows) == 0 {
				continue
			}
			bestFinish = append(bestFinish, followingThrows...)
			bestFinish = append(bestFinish, possibleThrow)
			break
		}
	}

	if len(bestFinish) == 0 {
		return nil
	}

	return bestFinish
}

func computePossibleNextThrows(score, throwsLeft int, followingThrows []models.ThrowType) []models.ThrowType {
	possibleThrows := models.GetAllThrowTypes(true, false, false)
	for _, possibleThrow := range possibleThrows {
		if score-possibleThrow.ToPoints() == 0 {
			return append(followingThrows, possibleThrow)
		}
		if score-possibleThrow.ToPoints() > 0 && throwsLeft-1 > 0 {
			nextFollowingThrows := computePossibleNextThrows(score-possibleThrow.ToPoints(), throwsLeft-1, followingThrows)
			if len(nextFollowingThrows) == 0 {
				continue
			}

			nextFollowingThrows = append(nextFollowingThrows, possibleThrow)
			return append(followingThrows, nextFollowingThrows...)
		}
	}

	return nil
}
