package actions_test

import (
	"encoding/json"
	"regexp"
	"server/actions"
	"server/mock_roller"
	"testing"
)

func logGameManager(t *testing.T, gameManager *actions.GameManager) {
	jsonBytes, _ := json.Marshal(gameManager)
	t.Logf(string(jsonBytes))
}

func TestNewGameManager(t *testing.T) {
	gameManager := actions.NewGameManager("http://localhost")
	if gameManager.ActiveActions.NewGameAction == nil {
		t.Fatalf(`expected non-nil new game action but got nil`)
	}
	if gameManager.ActiveActions.RollAction != nil {
		t.Fatalf(`expected nil roll action but got non-nil`)
	}
	if gameManager.ActiveActions.NewTurnAction != nil {
		t.Fatalf(`expected nil turn action but got non-nil`)
	}
	action, ok := gameManager.ActionLinks[actions.NewGameAction]
	if !ok {
		t.Fatalf(`expected a new game action but none were in the map`)
	}

	if action.Method != `POST` {
		t.Fatalf(`expected POST but got %v`, action.Method)
	}
	expectedUrl := `http://localhost/NewGame\?token=[0-9]+&winningScore=\{WinningScore\}`
	if matches, _ := regexp.MatchString(expectedUrl, action.Url); !matches {
		t.Fatalf(`expected '%v' but got '%v'`, expectedUrl, action.Url)
	}
	t.Logf(`url: %v`, action.Url)
	logGameManager(t, gameManager)
}

func TestGenerateNewGame(t *testing.T) {
	gameManager := actions.NewGameManager(`http://localhost`)
	gameManager.ActiveActions.NewGameAction(10000)
	if gameManager.ActiveGame == nil {
		t.Fatalf(`expected non-nil game but got nil`)
	}
	if gameManager.ActiveActions.NewGameAction == nil {
		t.Fatalf(`expected non-nil new game action but got nil`)
	}
	if gameManager.ActiveActions.NewTurnAction != nil {
		t.Fatalf(`expected nil turn action but got non-nil`)
	}
	if gameManager.ActiveActions.RollAction == nil {
		t.Fatalf(`expected non-nil roll action but got nil`)
	}

	action, ok := gameManager.ActionLinks[actions.RollAction]
	if !ok {
		t.Fatalf(`expected a roll action but none were in the map`)
	}

	if action.Method != `POST` {
		t.Fatalf(`expected POST but got %v`, action.Method)
	}
	expectedUrl := `http://localhost/Roll\?token=[0-9]+`
	if matches, _ := regexp.MatchString(expectedUrl, action.Url); !matches {
		t.Fatalf(`expected '%v' but got '%v'`, expectedUrl, action.Url)
	}
	t.Logf(`url: %v`, action.Url)
	logGameManager(t, gameManager)
}

func TestRollGame(t *testing.T) {
	gameManager := actions.NewGameManager(`http://localhost`)
	gameManager.ActiveActions.NewGameAction(10000)
	mockRoll := mock_roller.MockRoller{
		RollOverride: []int{1, 2, 2, 3, 3, 4},
	}
	gameManager.ActiveGame.Roller = mockRoll.Roll
	gameManager.ActiveActions.RollAction()

	if gameManager.ActiveGame == nil {
		t.Fatalf(`expected non-nil game but got nil`)
	}
	if gameManager.ActiveActions.NewGameAction == nil {
		t.Fatalf(`expected non-nil new game action but got nil`)
	}
	if gameManager.ActiveActions.NewTurnAction == nil {
		t.Fatalf(`expected non-nil turn action but got nil`)
	}
	if gameManager.ActiveActions.RollAction == nil {
		t.Fatalf(`expected non-nil roll action but got nil`)
	}

	action, ok := gameManager.ActionLinks[actions.NewTurnAction]
	if !ok {
		t.Fatalf(`expected a new turn action but none were in the map`)
	}

	if action.Method != `POST` {
		t.Fatalf(`expected POST but got %v`, action.Method)
	}
	expectedUrl := `http://localhost/NewTurn\?token=[0-9]+`
	if matches, _ := regexp.MatchString(expectedUrl, action.Url); !matches {
		t.Fatalf(`expected '%v' but got '%v'`, expectedUrl, action.Url)
	}
	t.Logf(`url: %v`, action.Url)
	logGameManager(t, gameManager)
}

func TestNoRollActionIfLastRollIsZero(t *testing.T) {
	gameManager := actions.NewGameManager(`http://localhost`)
	gameManager.ActiveActions.NewGameAction(10000)
	mockRoll := mock_roller.MockRoller{
		RollOverride: []int{2},
	}
	gameManager.ActiveGame.Roller = mockRoll.Roll
	gameManager.ActiveActions.RollAction()

	if gameManager.ActiveActions.NewTurnAction == nil {
		t.Fatalf(`expected non-nil turn action but got nil`)
	}
	if gameManager.ActiveActions.RollAction != nil {
		t.Fatalf(`expected nil roll action but got non-nil`)
	}
}

func TestRollAndStartNewTurnScore(t *testing.T) {
	gameManager := actions.NewGameManager(`http://localhost`)
	gameManager.ActiveActions.NewGameAction(10000)
	mockRoll := mock_roller.MockRoller{
		RollOverride: []int{1, 2, 2, 3, 3, 4},
	}
	gameManager.ActiveGame.Roller = mockRoll.Roll
	gameManager.ActiveActions.RollAction()
	gameManager.ActiveActions.NewTurnAction()

	if gameManager.ActiveActions.NewTurnAction != nil {
		t.Fatalf(`expected nil turn action but got non-nil`)
	}
	if gameManager.ActiveActions.RollAction == nil {
		t.Fatalf(`expected non-nil roll action but got non-nil`)
	}
	if gameManager.ActiveGame.CurrentScore != 100 {
		t.Fatalf(`expected score of 100 but got %v`, gameManager.ActiveGame.CurrentScore)
	}
	if totalTurns := len(gameManager.ActiveGame.Turns); totalTurns != 1 {
		t.Fatalf(`expected 1 turn but got %v`, totalTurns)
	}
}

func TestWinGame(t *testing.T) {
	gameManager := actions.NewGameManager(`http://localhost`)
	gameManager.ActiveActions.NewGameAction(10000)
	mockRoll := mock_roller.MockRoller{
		RollOverride: []int{1, 1, 1, 1, 1, 1},
	}
	gameManager.ActiveGame.Roller = mockRoll.Roll
	gameManager.ActiveActions.RollAction()

	logGameManager(t, gameManager)

	if gameManager.ActiveActions.NewGameAction == nil {
		t.Fatalf(`expected non-nil new game action but got nil`)
	}
	if gameManager.ActiveActions.NewTurnAction != nil {
		t.Fatalf(`expected nil turn action but got non-nil`)
	}
	if gameManager.ActiveActions.RollAction != nil {
		t.Fatalf(`expected nil roll action but got non-nil`)
	}
}
