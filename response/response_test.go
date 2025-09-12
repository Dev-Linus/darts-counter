package response

import (
	"darts-counter/models"
	"testing"
)

// helper to create a match with scores
func newMatchWithScore(score int, endMode models.IO, currentThrow uint32) *models.Match {
	return &models.Match{
		ID:            "test-match",
		Players:       []string{"p1"},
		CurrentPlayer: "p1",
		Scores:        map[string]int{"p1": score},
		CurrentThrow:  currentThrow,                  // 0, 1, 2, or 3 depending on turn progress
		EndMode:       models.MapIOToNumber(endMode), // 'd' for double-out etc.
	}
}

func TestGetPossibleFinishForMatchPlayer_ExactDoubleOut(t *testing.T) {
	match := newMatchWithScore(40, models.Double, 0)

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) == 0 {
		t.Fatalf("expected a finish but got none")
	}
	if finishes[len(finishes)-1] != models.D20 {
		t.Errorf("expected D20 as last dart, got %v", finishes[len(finishes)-1])
	}
}

func TestGetPossibleFinishForMatchPlayer_NoFinishTooHigh(t *testing.T) {
	match := newMatchWithScore(200, models.Double, 0)

	finishes := getPossibleFinishForMatchPlayer(match)

	if finishes != nil {
		t.Errorf("expected nil, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_CheckoutsInThree(t *testing.T) {
	match := newMatchWithScore(100, models.Double, 0)

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) < 2 {
		t.Fatalf("expected a 2-3 dart finish, got %v", finishes)
	}
	last := finishes[len(finishes)-1]
	if !last.IsDouble() {
		t.Errorf("expected last dart to be a double, got %v", last)
	}
}

func TestGetPossibleFinishForMatchPlayer_WhenTwoThrowsAlreadyUsed(t *testing.T) {
	// Only 1 dart left, player has 50 points, should only be possible with Bull
	match := newMatchWithScore(50, models.Straight, 2) // 2 throws used, only 1 left

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 1 || finishes[0] != models.BULL {
		t.Errorf("expected Bull finish, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_WhenTwoThrowsAlreadyUsed_MasterOut(t *testing.T) {
	// Only 1 dart left, player has 50 points, should only be possible with Bull
	match := newMatchWithScore(50, models.Master, 2) // 2 throws used, only 1 left

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 1 && finishes[0] == models.BULL {
		t.Errorf("expected nil finish, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_WhenThreeThrowsUsed_NoThrowsLeft(t *testing.T) {
	match := newMatchWithScore(32, models.Double, 3) // already 3 darts thrown

	finishes := getPossibleFinishForMatchPlayer(match)

	if finishes != nil {
		t.Errorf("expected nil because no darts left, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_MaximumScore180(t *testing.T) {
	// Player has 180 points left, only possible with 3x T20 (but not a finishing throw since endMode requires double)
	// Here we test that the calculation does not incorrectly return a "finish"
	match := newMatchWithScore(180, models.Double, 0)

	finishes := getPossibleFinishForMatchPlayer(match)

	if finishes != nil {
		t.Errorf("expected nil because 180 is not a finish with double-out, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_OneDartFinishOnDouble(t *testing.T) {
	match := newMatchWithScore(32, models.Double, 0) // Classic D16 finish

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 1 || finishes[0] != models.D16 {
		t.Errorf("expected D16 finish, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_HighestFinish(t *testing.T) {
	match := newMatchWithScore(180, models.Master, 0) // T20 T20 T20

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 3 || finishes[0] != models.T20 || finishes[1] != models.T20 || finishes[2] != models.T20 {
		t.Errorf("expected T20, T20, T20 finish, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_RandomStraightOut(t *testing.T) {
	match := newMatchWithScore(164, models.Straight, 0) // just random Straight out

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 3 || finishes[0] != models.BULL || finishes[1] != models.T18 || finishes[2] != models.T20 {
		t.Errorf("expected T20, T20, T20 finish, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_RandomDoubleOut(t *testing.T) {
	match := newMatchWithScore(164, models.Double, 0) // just random double out

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 3 || finishes[0] != models.T18 || finishes[1] != models.T20 || finishes[2] != models.BULL {
		t.Errorf("expected T18, T20, BULL finish, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_RandomMasterOut(t *testing.T) {
	match := newMatchWithScore(164, models.Master, 0) // just random master out

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 3 || finishes[0] != models.BULL || finishes[1] != models.T18 || finishes[2] != models.T20 {
		t.Errorf("expected T20, T20, T20 finish, got %v", finishes)
	}
}

func TestGetPossibleFinishForMatchPlayer_HighestDoubleOut(t *testing.T) {
	match := newMatchWithScore(170, models.Double, 0) // highest double out

	finishes := getPossibleFinishForMatchPlayer(match)

	if len(finishes) != 3 || finishes[0] != models.T20 || finishes[1] != models.T20 || finishes[2] != models.BULL {
		t.Errorf("expected T20, T20, BULL finish, got %v", finishes)
	}
}
