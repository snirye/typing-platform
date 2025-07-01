package core

import (
	"math"
	"time"
)

// NewGame creates a new game instance
func NewGame() *Game {
	game := &Game{
		State:       StateMenu,
		ScrollSpeed: 10.0, // pixels per second
		WordManager: NewWordManager(),
		ShouldExit:  false,
	}

	return game
}

// Start initializes the game with given dimensions
func (g *Game) Start(width, height int) {
	g.Width = width
	g.Height = height
	g.Renderer = NewRenderer(width, height)
	g.reset()
}

// UpdateDimensions updates the game dimensions
func (g *Game) UpdateDimensions(width, height int) {
	g.Width = width
	g.Height = height
	if g.Renderer != nil {
		g.Renderer.UpdateDimensions(width, height)
	}
}

// ProcessInput handles user input
func (g *Game) ProcessInput(key rune) {
	switch g.State {
	case StateMenu:
		g.processMenuInput(key)
	case StatePlaying:
		g.processGameInput(key)
	case StatePaused:
		g.processPauseInput(key)
	case StateGameOver:
		g.processGameOverInput(key)
	}
}

// Render updates game logic and returns the rendered frame
func (g *Game) Render() string {
	if g.State == StatePlaying {
		g.updateGameLogic()
	}

	if g.Renderer != nil {
		return g.Renderer.RenderGame(g)
	}
	return "Renderer not initialized"
}

// ShouldQuit returns whether the game should exit
func (g *Game) ShouldQuit() bool {
	return g.ShouldExit
}

// Private methods

func (g *Game) reset() {
	g.Score = 0
	g.WordsTyped = 0
	g.CharsTyped = 0
	g.StartTime = time.Now()
	g.ScrollOffset = 0
	g.Player = Player{
		X:        g.Width / 2,
		Y:        g.Height - 10,
		Platform: 0,
	}

	// Initialize platforms
	g.generateInitialPlatforms()
}

func (g *Game) processMenuInput(key rune) {
	switch key {
	case ' ': // Space to start
		g.State = StatePlaying
		g.reset()
	case 'q', 'Q':
		g.ShouldExit = true
	case 27: // ESC
		g.ShouldExit = true
	}
}

func (g *Game) processGameInput(key rune) {
	switch key {
	case 27: // ESC - pause
		g.State = StatePaused
	case 8, 127: // Backspace
		g.handleBackspace()
	default:
		if isAlphanumeric(key) {
			g.handleTyping(key)
		}
	}
}

func (g *Game) processPauseInput(key rune) {
	switch key {
	case 27: // ESC - resume
		g.State = StatePlaying
	case 'q', 'Q':
		g.ShouldExit = true
	}
}

func (g *Game) processGameOverInput(key rune) {
	switch key {
	case ' ': // Space to restart
		g.State = StatePlaying
		g.reset()
	case 'q', 'Q':
		g.ShouldExit = true
	}
}

func (g *Game) handleTyping(key rune) {
	if len(g.Platforms) == 0 {
		return
	}

	currentPlatform := &g.Platforms[g.Player.Platform]

	// Check if the character is correct
	if g.WordManager.IsValidChar(currentPlatform.Word, currentPlatform.Typed, key) {
		currentPlatform.Typed += string(key)
		g.CharsTyped++

		// Check if word is complete
		if g.WordManager.IsWordComplete(currentPlatform.Word, currentPlatform.Typed) {
			g.completeWord(currentPlatform)
		}
	}
}

func (g *Game) handleBackspace() {
	if len(g.Platforms) == 0 {
		return
	}

	currentPlatform := &g.Platforms[g.Player.Platform]
	if len(currentPlatform.Typed) > 0 {
		currentPlatform.Typed = currentPlatform.Typed[:len(currentPlatform.Typed)-1]
	}
}

func (g *Game) completeWord(platform *Platform) {
	platform.Complete = true
	g.WordsTyped++
	g.Score += len(platform.Word) * 10 // Base score

	// Bonus for speed - reward faster typing
	timeSinceStart := time.Since(g.StartTime).Seconds()
	if timeSinceStart > 0 {
		speedBonus := int(math.Max(0, 100-(timeSinceStart/float64(g.WordsTyped))))
		g.Score += speedBonus
	}

	// Move player to next platform
	g.jumpToNextPlatform()
}

func (g *Game) jumpToNextPlatform() {
	// Find next available platform above current one
	currentY := g.Platforms[g.Player.Platform].Y
	nextPlatformIndex := -1

	for i, platform := range g.Platforms {
		if platform.Y < currentY && !platform.Complete {
			if nextPlatformIndex == -1 || platform.Y > g.Platforms[nextPlatformIndex].Y {
				nextPlatformIndex = i
			}
		}
	}

	if nextPlatformIndex != -1 {
		g.Player.Platform = nextPlatformIndex
		platform := g.Platforms[nextPlatformIndex]
		g.Player.X = platform.X + platform.Width/2
		g.Player.Y = platform.Y - 1
	} else {
		// Generate new platforms if needed
		g.generateMorePlatforms()
	}
}

func (g *Game) updateGameLogic() {
	// Update scroll offset
	deltaTime := 1.0 / 60.0 // Assume 60 FPS
	g.ScrollOffset += g.ScrollSpeed * deltaTime

	// Update player position based on scroll
	g.Player.Y = g.Platforms[g.Player.Platform].Y - 1 - int(g.ScrollOffset)

	// Check if player fell off screen
	if g.Player.Y >= g.Height-3 {
		g.State = StateGameOver
		return
	}

	// Generate new platforms as needed
	g.generateMorePlatforms()

	// Remove old platforms that are off screen
	g.cleanupPlatforms()
}

func (g *Game) generateInitialPlatforms() {
	g.Platforms = make([]Platform, 0)

	// Generate starting platform
	startPlatform := Platform{
		X:        g.Width/2 - 10,
		Y:        g.Height - 5,
		Width:    20,
		Word:     g.WordManager.GetRandomWord(),
		Typed:    "",
		Complete: false,
	}
	g.Platforms = append(g.Platforms, startPlatform)

	// Generate additional platforms
	for i := 1; i < 10; i++ {
		platform := Platform{
			X:        g.Width/4 + (i%2)*(g.Width/2),
			Y:        g.Height - 5 - i*50,
			Width:    15 + (i%3)*10,
			Word:     g.WordManager.GetRandomWord(),
			Typed:    "",
			Complete: false,
		}
		g.Platforms = append(g.Platforms, platform)
	}
}

func (g *Game) generateMorePlatforms() {
	if len(g.Platforms) == 0 {
		return
	}

	// Find the highest platform
	highestY := g.Platforms[0].Y
	for _, platform := range g.Platforms {
		if platform.Y < highestY {
			highestY = platform.Y
		}
	}

	// Generate new platforms above the highest one
	topOfScreen := int(g.ScrollOffset) - 100
	if highestY > topOfScreen {
		numNewPlatforms := 5
		for i := 0; i < numNewPlatforms; i++ {
			platform := Platform{
				X:        20 + (i%3)*((g.Width-40)/3),
				Y:        highestY - 60 - i*50,
				Width:    10 + (i%4)*5,
				Word:     g.WordManager.GetRandomWord(),
				Typed:    "",
				Complete: false,
			}
			g.Platforms = append(g.Platforms, platform)
		}
	}
}

func (g *Game) cleanupPlatforms() {
	// Remove platforms that are far below the screen
	bottomThreshold := int(g.ScrollOffset) + g.Height + 100

	newPlatforms := make([]Platform, 0)
	for i, platform := range g.Platforms {
		if platform.Y < bottomThreshold {
			newPlatforms = append(newPlatforms, platform)
		} else if i == g.Player.Platform {
			// Don't remove the platform the player is on
			newPlatforms = append(newPlatforms, platform)
		}
	}

	// Update player platform index if needed
	if len(newPlatforms) < len(g.Platforms) {
		for i, platform := range newPlatforms {
			if platform.X == g.Platforms[g.Player.Platform].X &&
				platform.Y == g.Platforms[g.Player.Platform].Y {
				g.Player.Platform = i
				break
			}
		}
	}

	g.Platforms = newPlatforms
}

// GetStats returns current game statistics
func (g *Game) GetStats() Stats {
	gameTime := time.Since(g.StartTime)
	minutes := gameTime.Minutes()

	wpm := 0.0
	cpm := 0.0

	if minutes > 0 {
		wpm = float64(g.WordsTyped) / minutes
		cpm = float64(g.CharsTyped) / minutes
	}

	return Stats{
		Score:      g.Score,
		WPM:        wpm,
		CPM:        cpm,
		WordsTyped: g.WordsTyped,
		CharsTyped: g.CharsTyped,
		GameTime:   gameTime,
	}
}

// Helper functions
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
