package core

import (
	"testing"
	"time"
)

func TestNewGame(t *testing.T) {
	game, err := NewGame("test_log.txt")
	if err != nil {
		t.Fatalf("NewGame() error: %v", err)
	}

	if game == nil {
		t.Fatal("NewGame() returned nil")
	}

	if game.State != StateMenu {
		t.Errorf("Expected initial state to be StateMenu, got %v", game.State)
	}

	if game.WordManager == nil {
		t.Error("WordManager should be initialized")
	}
}

func TestGameStart(t *testing.T) {
	game, err := NewGame("test_log.txt")
	if err != nil {
		t.Fatalf("NewGame() error: %v", err)
	}
	game.Start(80, 24)

	if game.Width != 80 || game.Height != 24 {
		t.Errorf("Expected dimensions 80x24, got %dx%d", game.Width, game.Height)
	}

	if game.Renderer == nil {
		t.Error("Renderer should be initialized after Start()")
	}

	if len(game.Platforms) == 0 {
		t.Error("Platforms should be generated after Start()")
	}
}

func TestProcessMenuInput(t *testing.T) {
	game, err := NewGame("test_log.txt")
	if err != nil {
		t.Fatalf("NewGame() error: %v", err)
	}
	game.Start(80, 24)

	// Test space key starts the game
	game.ProcessInput(' ')
	if game.State != StatePlaying {
		t.Errorf("Expected state to be StatePlaying after space, got %v", game.State)
	}

	// Reset to menu
	game.State = StateMenu

	// Test quit key
	game.ProcessInput('q')
	if !game.ShouldQuit() {
		t.Error("Expected ShouldQuit to be true after 'q'")
	}
}

func TestTypingValidation(t *testing.T) {
	game, err := NewGame("test_log.txt")
	if err != nil {
		t.Fatalf("NewGame() error: %v", err)
	}
	game.Start(80, 24)
	game.State = StatePlaying

	// Ensure we have a platform with a word
	if len(game.Platforms) == 0 {
		t.Fatal("No platforms generated")
	}

	platform := &game.Platforms[game.Player.Platform]
	originalWord := platform.Word

	// Test valid character
	if len(originalWord) > 0 {
		firstChar := rune(originalWord[0])
		game.ProcessInput(firstChar)

		if len(platform.Typed) != 1 {
			t.Errorf("Expected typed length 1, got %d", len(platform.Typed))
		}

		if platform.Typed != string(firstChar) {
			t.Errorf("Expected typed '%s', got '%s'", string(firstChar), platform.Typed)
		}
	}
}

func TestBackspace(t *testing.T) {
	game, err := NewGame("test_log.txt")
	if err != nil {
		t.Fatalf("NewGame() error: %v", err)
	}
	game.Start(80, 24)
	game.State = StatePlaying

	platform := &game.Platforms[game.Player.Platform]
	platform.Typed = "test"

	// Test backspace
	game.ProcessInput(8) // Backspace

	if platform.Typed != "tes" {
		t.Errorf("Expected 'tes' after backspace, got '%s'", platform.Typed)
	}
}

func TestStats(t *testing.T) {
	game, err := NewGame("test_log.txt")
	if err != nil {
		t.Fatalf("NewGame() error: %v", err)
	}
	game.Start(80, 24)

	// Simulate some typing
	game.WordsTyped = 5
	game.CharsTyped = 25
	game.StartTime = time.Now().Add(-1 * time.Minute) // 1 minute ago

	stats := game.GetStats()

	if stats.WordsTyped != 5 {
		t.Errorf("Expected 5 words typed, got %d", stats.WordsTyped)
	}

	if stats.CharsTyped != 25 {
		t.Errorf("Expected 25 chars typed, got %d", stats.CharsTyped)
	}

	// WPM should be approximately 5 (5 words in 1 minute)
	if stats.WPM < 4 || stats.WPM > 6 {
		t.Errorf("Expected WPM around 5, got %.2f", stats.WPM)
	}
}
