package core

import (
	"math"
	"time"
)

// NewGame creates a new game instance with logging to the specified file
// logsPath: path to the log file for debug output
func NewGame(logsPath string) (*Game, error) {
	logger, err := NewLogger(logsPath)
	if err != nil {
		return nil, err
	}
	logger.Println("NewGame: initializing game")
	game := &Game{
		State:       StateMenu,
		ScrollSpeed: 5.0, // pixels per second - increased for visible scrolling. default to 5.0
		WordManager: NewWordManager(),
		ShouldExit:  false,
		Logger:      logger,
	}
	logger.Println("NewGame: game struct created")
	return game, nil
}

// Start initializes the game with given dimensions
func (g *Game) Start(width, height int) {
	g.Logger.Printf("Start: width=%d, height=%d", width, height)
	g.Width = width
	g.Height = height
	g.Renderer = NewRenderer(width, height)
	g.reset()
}

// UpdateDimensions updates the game dimensions
func (g *Game) UpdateDimensions(width, height int) {
	g.Logger.Printf("UpdateDimensions: width=%d, height=%d", width, height)
	g.Width = width
	g.Height = height
	if g.Renderer != nil {
		g.Renderer.UpdateDimensions(width, height)
	}
}

// ProcessInput handles user input
func (g *Game) ProcessInput(key rune) {
	g.Logger.Printf("ProcessInput: key=%v, state=%v", key, g.State)
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
	// g.Logger.Printf("Render: state=%v", g.State)
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
	g.Logger.Println("reset: resetting game state")
	g.Score = 0
	g.WordsTyped = 0
	g.CharsTyped = 0
	g.StartTime = time.Now()
	g.ScrollOffset = 0
	g.ScrollAccumulator = 0 // Reset scroll accumulator
	g.Player = Player{
		X:        g.Width / 2,
		Y:        g.Height/4 - 1, // Position player on the starting platform in upper portion
		Platform: 0,
	}

	// Initialize platforms
	g.generateInitialPlatforms()
}

func (g *Game) processMenuInput(key rune) {
	g.Logger.Printf("processMenuInput: key=%v", key)
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
	g.Logger.Printf("processGameInput: key=%v", key)
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
	g.Logger.Printf("processPauseInput: key=%v", key)
	switch key {
	case 27: // ESC - resume
		g.State = StatePlaying
	case 'q', 'Q':
		g.ShouldExit = true
	}
}

func (g *Game) processGameOverInput(key rune) {
	g.Logger.Printf("processGameOverInput: key=%v", key)
	switch key {
	case ' ': // Space to restart
		g.State = StatePlaying
		g.reset()
	case 'q', 'Q':
		g.ShouldExit = true
	}
}

func (g *Game) handleTyping(key rune) {
	g.Logger.Printf("handleTyping: key=%v", key)
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
	g.Logger.Println("handleBackspace")
	if len(g.Platforms) == 0 {
		return
	}

	currentPlatform := &g.Platforms[g.Player.Platform]
	if len(currentPlatform.Typed) > 0 {
		currentPlatform.Typed = currentPlatform.Typed[:len(currentPlatform.Typed)-1]
	}
}

func (g *Game) completeWord(platform *Platform) {
	g.Logger.Printf("completeWord: word=%s", platform.Word)
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
	g.Logger.Println("jumpToNextPlatform")
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
	// g.Logger.Println("updateGameLogic")
	// Calculate scroll movement - platforms scroll down, creating upward movement effect
	deltaTime := 1.0 / 60.0 // Assume 60 FPS
	scrollDelta := g.ScrollSpeed * deltaTime

	// Accumulate fractional scroll amounts - this ensures smooth scrolling at any speed
	g.ScrollAccumulator += scrollDelta

	// Only move platforms when we have accumulated at least 1 full pixel
	pixelMovement := 0
	if g.ScrollAccumulator >= 1.0 {
		pixelMovement = int(g.ScrollAccumulator)
		g.ScrollAccumulator -= float64(pixelMovement) // Keep the fractional remainder
	}

	// Debug: Log the scroll values
	g.Logger.Printf("updateGameLogic: scrollDelta=%.3f, accumulator=%.3f, pixelMovement=%d",
		scrollDelta, g.ScrollAccumulator, pixelMovement)

	// Update all platform positions only if we have movement
	if pixelMovement > 0 {
		for i := range g.Platforms {
			oldY := g.Platforms[i].Y
			g.Platforms[i].Y += pixelMovement
			if i == 0 { // Log first platform movement for debugging
				g.Logger.Printf("Platform 0: Y changed from %d to %d (moved %d pixels)", oldY, g.Platforms[i].Y, pixelMovement)
			}
		}
	}

	// Update player position to match current platform movement
	if len(g.Platforms) > 0 && g.Player.Platform < len(g.Platforms) {
		platform := g.Platforms[g.Player.Platform]
		g.Player.X = platform.X + platform.Width/2
		g.Player.Y = platform.Y - 1 // Player sits on top of platform
	}

	// Check if player fell off screen (now using direct Y position)
	if g.Player.Y >= g.Height-3 {
		g.Logger.Println("updateGameLogic: player fell off screen, game over")
		g.State = StateGameOver
		return
	}

	// Generate new platforms as screen scrolls up
	g.generateMorePlatforms()

	// Remove old platforms that are off screen
	g.cleanupPlatforms()
}

func (g *Game) generateInitialPlatforms() {
	g.Logger.Println("generateInitialPlatforms")
	g.Platforms = make([]Platform, 0)

	// Generate starting platform near the top but with room for upward progression
	startPlatform := Platform{
		X:        g.Width/2 - 10,
		Y:        g.Height / 4, // Start in upper portion of screen
		Width:    20,
		Word:     g.WordManager.GetRandomWord(),
		Typed:    "",
		Complete: false,
	}
	g.Platforms = append(g.Platforms, startPlatform)

	// Generate platforms going upward (decreasing Y values) for progression
	currentY := startPlatform.Y
	for i := 1; i < 15; i++ { // Generate more initial platforms
		// Vary X position across screen width
		xPos := 20 + (i%4)*(g.Width-40)/4
		// Ensure minimum platform spacing going upward
		currentY -= 10 + (i % 3) // Vary vertical spacing upward

		platform := Platform{
			X:        xPos,
			Y:        currentY,
			Width:    15 + (i%3)*10,
			Word:     g.WordManager.GetRandomWord(),
			Typed:    "",
			Complete: false,
		}
		g.Platforms = append(g.Platforms, platform)
	}
}

func (g *Game) generateMorePlatforms() {
	// g.Logger.Println("generateMorePlatforms")
	if len(g.Platforms) == 0 {
		return
	}

	// Find the highest platform (lowest Y value)
	highestY := g.Platforms[0].Y
	for _, platform := range g.Platforms {
		if platform.Y < highestY {
			highestY = platform.Y
		}
	}

	// Generate new platforms when the highest platform gets close to being visible
	if highestY > -200 { // Generate when platforms are 200 pixels above screen
		numNewPlatforms := 8
		for i := 0; i < numNewPlatforms; i++ {
			// Vary X position across the screen
			xPos := 30 + (i%5)*(g.Width-60)/5
			// Place new platforms above the current highest
			newY := highestY - 60 - i*45 // Consistent upward spacing

			platform := Platform{
				X:        xPos,
				Y:        newY,
				Width:    12 + (i%4)*6,
				Word:     g.WordManager.GetRandomWord(),
				Typed:    "",
				Complete: false,
			}
			g.Platforms = append(g.Platforms, platform)
		}
	}
}

func (g *Game) cleanupPlatforms() {
	// g.Logger.Println("cleanupPlatforms")
	// Remove platforms that have scrolled far below the screen
	bottomThreshold := g.Height + 100

	newPlatforms := make([]Platform, 0)
	playerPlatformFound := false
	newPlayerPlatform := 0

	for i, platform := range g.Platforms {
		// Keep platforms that are still relevant (not too far below screen)
		if platform.Y < bottomThreshold {
			newPlatforms = append(newPlatforms, platform)
			// Track player's platform in the new array
			if i == g.Player.Platform {
				newPlayerPlatform = len(newPlatforms) - 1
				playerPlatformFound = true
			}
		}
	}

	// Update platforms array and player platform index
	g.Platforms = newPlatforms
	if playerPlatformFound {
		g.Player.Platform = newPlayerPlatform
	} else if len(g.Platforms) > 0 {
		// If player's platform was removed, move to closest available platform
		g.Player.Platform = 0
	}
}

// GetStats returns current game statistics
func (g *Game) GetStats() Stats {
	// g.Logger.Println("GetStats")
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
