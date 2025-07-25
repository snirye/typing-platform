package core

import (
	"fmt"
	"strings"
	"time"
)

// Colors for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Renderer handles ASCII art rendering for the game
type Renderer struct {
	width  int
	height int
}

// NewRenderer creates a new renderer
func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		width:  width,
		height: height,
	}
}

// UpdateDimensions updates the renderer dimensions
func (r *Renderer) UpdateDimensions(width, height int) {
	r.width = width
	r.height = height
}

// RenderGame renders the complete game state
func (r *Renderer) RenderGame(g *Game) string {
	switch g.State {
	case StateMenu:
		return r.renderMenu(g)
	case StatePlaying:
		return r.renderGameplay(g)
	case StatePaused:
		return r.renderPaused(g)
	case StateGameOver:
		return r.renderGameOver(g)
	default:
		return "Unknown game state"
	}
}

// renderMenu renders the main menu
func (r *Renderer) renderMenu(g *Game) string {
	var sb strings.Builder

	// Clear screen and position cursor
	sb.WriteString("\033[2J\033[H")

	// Center the menu
	centerY := r.height / 2
	centerX := r.width / 2

	// Title
	title := "ASCII TYPING PLATFORMER"
	titleX := centerX - len(title)/2
	r.writeAtPosition(&sb, titleX, centerY-3, ColorBold+ColorCyan+title+ColorReset)

	// Menu options
	options := []string{
		"Press SPACE to Start",
		"Press Q to Quit",
	}

	for i, option := range options {
		optionX := centerX - len(option)/2
		r.writeAtPosition(&sb, optionX, centerY+i, ColorWhite+option+ColorReset)
	}

	return sb.String()
}

// renderGameplay renders the main game view
func (r *Renderer) renderGameplay(g *Game) string {
	var sb strings.Builder

	// Use the logger for debug output. This shows actual platform positions since they're now updated directly by game logic.
	// for i, p := range g.Platforms {
	// 	g.Logger.Printf("Platform %d: X=%d Y=%d Complete=%v\n", i, p.X, p.Y, p.Complete)
	// }

	// Clear screen
	sb.WriteString("\033[2J\033[H")

	// Create a grid to draw on
	grid := make([][]rune, r.height)
	for i := range grid {
		grid[i] = make([]rune, r.width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Draw platforms
	for _, platform := range g.Platforms {
		r.drawPlatform(grid, platform)
	}

	// Draw player
	r.drawPlayer(grid, g.Player)

	// Convert grid to string
	for y := 0; y < r.height-4; y++ { // Leave space for HUD (4 lines: border + stats + wpm + current word)
		for x := 0; x < r.width; x++ {
			if x < len(grid[y]) && y < len(grid) {
				sb.WriteRune(grid[y][x])
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}

	// Draw HUD
	sb.WriteString(r.renderHUD(g))

	return sb.String()
}

// renderPaused renders the pause screen
func (r *Renderer) renderPaused(g *Game) string {
	var sb strings.Builder

	// Show game underneath
	sb.WriteString(r.renderGameplay(g))

	// Overlay pause message
	centerY := r.height / 2
	centerX := r.width / 2

	pauseMsg := "PAUSED"
	pauseX := centerX - len(pauseMsg)/2
	r.writeAtPosition(&sb, pauseX, centerY-1, ColorBold+ColorYellow+pauseMsg+ColorReset)

	resumeMsg := "Press ESC to resume, Q to quit"
	resumeX := centerX - len(resumeMsg)/2
	r.writeAtPosition(&sb, resumeX, centerY+1, ColorWhite+resumeMsg+ColorReset)

	return sb.String()
}

// renderGameOver renders the game over screen
func (r *Renderer) renderGameOver(g *Game) string {
	var sb strings.Builder

	// Clear screen
	sb.WriteString("\033[2J\033[H")

	centerY := r.height / 2
	centerX := r.width / 2

	// Game Over title
	gameOverMsg := "GAME OVER"
	gameOverX := centerX - len(gameOverMsg)/2
	r.writeAtPosition(&sb, gameOverX, centerY-4, ColorBold+ColorRed+gameOverMsg+ColorReset)

	// Stats
	stats := g.GetStats()
	statsLines := []string{
		fmt.Sprintf("Score: %d", stats.Score),
		fmt.Sprintf("WPM: %.1f", stats.WPM),
		fmt.Sprintf("CPM: %.1f", stats.CPM),
		fmt.Sprintf("Words: %d", stats.WordsTyped),
		fmt.Sprintf("Time: %s", formatDuration(stats.GameTime)),
	}

	for i, line := range statsLines {
		lineX := centerX - len(line)/2
		r.writeAtPosition(&sb, lineX, centerY-1+i, ColorWhite+line+ColorReset)
	}

	// Options
	optionsMsg := "Press SPACE to play again, Q to quit"
	optionsX := centerX - len(optionsMsg)/2
	r.writeAtPosition(&sb, optionsX, centerY+6, ColorGreen+optionsMsg+ColorReset)

	return sb.String()
}

// renderHUD renders the heads-up display
func (r *Renderer) renderHUD(g *Game) string {
	stats := g.GetStats()

	// Top border
	border := strings.Repeat("=", r.width)

	// HUD line 1: Score and time
	gameTime := time.Since(g.StartTime)
	line1 := fmt.Sprintf("Score: %d | Time: %s", stats.Score, formatDuration(gameTime))

	// HUD line 2: WPM and CPM
	line2 := fmt.Sprintf("WPM: %.1f | CPM: %.1f | Words: %d", stats.WPM, stats.CPM, stats.WordsTyped)

	// Current word display - always show status
	currentWord := ""
	if len(g.Platforms) > 0 && g.Player.Platform < len(g.Platforms) {
		platform := g.Platforms[g.Player.Platform]
		if !platform.Complete {
			// Show current word with progress highlighting
			typed := ColorGreen + "[" + platform.Typed + "]" + ColorReset
			remaining := ColorWhite + platform.Word[len(platform.Typed):] + ColorReset
			currentWord = fmt.Sprintf("Word: %s%s", typed, remaining)
		} else {
			// Show completed word in green
			currentWord = fmt.Sprintf("Word: %s%s%s (Complete!)", ColorGreen, platform.Word, ColorReset)
		}
	} else {
		// No active platform
		currentWord = "Word: " + ColorYellow + "No active word" + ColorReset
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n",
		border,
		r.padString(line1, r.width),
		r.padString(line2, r.width),
		r.padString(currentWord, r.width))
}

// Helper methods
func (r *Renderer) writeAtPosition(sb *strings.Builder, x, y int, text string) {
	if y >= 0 && y < r.height && x >= 0 {
		sb.WriteString(fmt.Sprintf("\033[%d;%dH%s", y+1, x+1, text))
	}
}

func (r *Renderer) drawPlatform(grid [][]rune, platform Platform) {
	// Platform Y position is now its actual screen position
	screenY := platform.Y

	// Only draw if platform is visible on screen
	if screenY >= 0 && screenY < len(grid)-3 {
		// Draw platform line
		for i := 0; i < platform.Width && platform.X+i < len(grid[screenY]); i++ {
			if platform.X+i >= 0 {
				grid[screenY][platform.X+i] = '='
			}
		}

		// Draw word below platform with typed indicator
		if screenY+1 < len(grid)-3 && !platform.Complete {
			// Create display string with typed characters in brackets
			typedPart := "[" + platform.Typed + "]"
			remainingPart := platform.Word[len(platform.Typed):]
			displayWord := typedPart + remainingPart

			wordX := platform.X + platform.Width/2 - len(displayWord)/2
			for i, char := range displayWord {
				if wordX+i >= 0 && wordX+i < len(grid[screenY+1]) {
					grid[screenY+1][wordX+i] = char
				}
			}
		}
	}
}

func (r *Renderer) drawPlayer(grid [][]rune, player Player) {
	// Player Y position is now its actual screen position
	screenY := player.Y

	if screenY >= 0 && screenY < len(grid)-3 && player.X >= 0 && player.X < len(grid[0]) {
		grid[screenY][player.X] = '@'
	}
}

func (r *Renderer) padString(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
