package core

import "time"

// GameInterface defines the public interface for the game engine
type GameInterface interface {
	Start(width, height int)
	UpdateDimensions(width, height int)
	ProcessInput(key rune) // No return value - game manages its own state
	Render() string        // Updates game logic and returns rendered frame
	ShouldQuit() bool      // indicates exit was performed
}

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// Player represents the player character
type Player struct {
	X, Y     int
	Platform int // Current platform index
}

// Platform represents a platform in the game
type Platform struct {
	X, Y     int
	Width    int
	Word     string
	Typed    string
	Complete bool
}

// Game holds the game state and logic
type Game struct {
	State             GameState
	Width             int
	Height            int
	Player            Player
	Platforms         []Platform
	Score             int
	StartTime         time.Time
	WordsTyped        int
	CharsTyped        int
	ShouldExit        bool
	ScrollSpeed       float64
	ScrollOffset      float64
	ScrollAccumulator float64 // Accumulates fractional scroll amounts
	WordManager       *WordManager
	Renderer          *Renderer
	Logger            *Logger // Add a Logger field for debug logging
}

// Stats represents game statistics
type Stats struct {
	Score      int
	WPM        float64
	CPM        float64
	Accuracy   float64
	WordsTyped int
	CharsTyped int
	GameTime   time.Duration
}
