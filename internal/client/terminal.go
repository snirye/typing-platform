package client

import (
	"ascii-type/internal/core"
	"time"

	"github.com/nsf/termbox-go"
)

// TerminalClient handles terminal I/O and display
type TerminalClient struct {
	game   core.GameInterface
	width  int
	height int
}

// NewTerminalClient creates a new terminal client
func NewTerminalClient(game core.GameInterface) *TerminalClient {
	return &TerminalClient{
		game: game,
	}
}

// Run starts the main game loop
func (tc *TerminalClient) Run() error {
	// Initialize termbox
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	// Set input and output modes
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.SetOutputMode(termbox.OutputNormal)

	// Get initial terminal size
	tc.width, tc.height = termbox.Size()
	tc.game.Start(tc.width, tc.height)

	// Create channels for events and ticker
	eventChan := make(chan termbox.Event)
	ticker := time.NewTicker(time.Second / 60) // 60 FPS
	defer ticker.Stop()

	// Start event polling goroutine
	go func() {
		for {
			event := termbox.PollEvent()
			eventChan <- event
		}
	}()

	// Main game loop
	for !tc.game.ShouldQuit() {
		select {
		case event := <-eventChan:
			if !tc.handleEvent(event) {
				return nil // Exit requested
			}

		case <-ticker.C:
			// Render game frame
			tc.render()
		}
	}

	return nil
}

// handleEvent processes termbox events
func (tc *TerminalClient) handleEvent(event termbox.Event) bool {
	switch event.Type {
	case termbox.EventKey:
		// Handle special keys
		if event.Key == termbox.KeyEsc {
			tc.game.ProcessInput(27) // ESC
		} else if event.Key == termbox.KeySpace {
			tc.game.ProcessInput(' ')
		} else if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 {
			tc.game.ProcessInput(8) // Backspace
		} else if event.Key == termbox.KeyCtrlC {
			return false // Exit
		} else if event.Ch != 0 {
			// Regular character
			tc.game.ProcessInput(event.Ch)
		}

	case termbox.EventResize:
		// Handle terminal resize
		tc.width, tc.height = termbox.Size()
		tc.game.UpdateDimensions(tc.width, tc.height)

	case termbox.EventError:
		return false
	}

	return true
}

// render clears the screen and draws the current game frame
func (tc *TerminalClient) render() {
	// Clear the screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Get rendered frame from game
	frame := tc.game.Render()

	// Draw frame to terminal
	tc.drawFrame(frame)

	// Flush to screen
	termbox.Flush()
}

// drawFrame renders the game frame to the terminal
func (tc *TerminalClient) drawFrame(frame string) {
	x, y := 0, 0

	// Parse ANSI escape sequences for basic color support
	inEscape := false
	escapeSeq := ""
	currentFg := termbox.ColorDefault
	currentBg := termbox.ColorDefault

	for _, ch := range frame {
		if ch == '\033' {
			inEscape = true
			escapeSeq = string(ch)
			continue
		}

		if inEscape {
			escapeSeq += string(ch)
			if ch == 'm' {
				// End of escape sequence
				inEscape = false
				currentFg, currentBg = tc.parseColor(escapeSeq)
				escapeSeq = ""
			} else if ch == 'H' {
				// Cursor position - extract coordinates if needed
				inEscape = false
				escapeSeq = ""
			} else if ch == 'J' {
				// Clear screen command
				inEscape = false
				escapeSeq = ""
			}
			continue
		}

		if ch == '\n' {
			x = 0
			y++
		} else if ch == '\r' {
			x = 0
		} else {
			if x < tc.width && y < tc.height {
				termbox.SetCell(x, y, ch, currentFg, currentBg)
			}
			x++
		}
	}
}

// parseColor converts ANSI color codes to termbox colors
func (tc *TerminalClient) parseColor(escapeSeq string) (termbox.Attribute, termbox.Attribute) {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	// Simple color mapping for basic colors
	switch escapeSeq {
	case "\033[0m": // Reset
		fg = termbox.ColorDefault
		bg = termbox.ColorDefault
	case "\033[31m": // Red
		fg = termbox.ColorRed
	case "\033[32m": // Green
		fg = termbox.ColorGreen
	case "\033[33m": // Yellow
		fg = termbox.ColorYellow
	case "\033[34m": // Blue
		fg = termbox.ColorBlue
	case "\033[35m": // Magenta
		fg = termbox.ColorMagenta
	case "\033[36m": // Cyan
		fg = termbox.ColorCyan
	case "\033[37m": // White
		fg = termbox.ColorWhite
	case "\033[1m": // Bold (bright)
		fg = fg | termbox.AttrBold
	}

	return fg, bg
}
