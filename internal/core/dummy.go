package core

import (
	"fmt"
	"strings"
	"time"
)

// DummyGame implements GameInterface for testing client implementations
// It logs method calls and displays them in the rendered output
type DummyGame struct {
	width       int
	height      int
	shouldQuit  bool
	messages    []string // Buffer to store method call messages
	maxMessages int      // Maximum number of messages to keep
}

// NewDummyGame creates a new dummy game instance
func NewDummyGame() *DummyGame {
	return &DummyGame{
		messages:    make([]string, 0),
		maxMessages: 10, // Keep last 10 messages to avoid clutter
	}
}

// addMessage adds a new message to the buffer, maintaining max size
func (d *DummyGame) addMessage(message string) {
	timestamp := time.Now().Format("15:04:05")
	fullMessage := fmt.Sprintf("[%s] %s", timestamp, message)

	d.messages = append(d.messages, fullMessage)

	// Keep only the last maxMessages to prevent overflow
	if len(d.messages) > d.maxMessages {
		d.messages = d.messages[len(d.messages)-d.maxMessages:]
	}
}

// Start initializes the dummy game with given dimensions
func (d *DummyGame) Start(width, height int) {
	d.width = width
	d.height = height
	d.shouldQuit = false
	d.addMessage(fmt.Sprintf("Start called with width=%d, height=%d", width, height))
}

// UpdateDimensions updates the game dimensions and logs the call
func (d *DummyGame) UpdateDimensions(width, height int) {
	d.width = width
	d.height = height
	d.addMessage(fmt.Sprintf("UpdateDimensions called with width=%d, height=%d", width, height))
}

// ProcessInput processes input and logs the call with the key received
func (d *DummyGame) ProcessInput(key rune) {
	// Handle quit command
	if key == 'q' || key == 'Q' {
		d.shouldQuit = true
		d.addMessage("ProcessInput called with key='q' - quit requested")
		return
	}

	// Log the input with readable representation
	var keyDesc string
	switch key {
	case '\n', '\r':
		keyDesc = "ENTER"
	case '\t':
		keyDesc = "TAB"
	case ' ':
		keyDesc = "SPACE"
	case 27: // ESC
		keyDesc = "ESC"
	default:
		if key >= 32 && key <= 126 { // Printable ASCII
			keyDesc = fmt.Sprintf("'%c'", key)
		} else {
			keyDesc = fmt.Sprintf("(code:%d)", key)
		}
	}

	d.addMessage(fmt.Sprintf("ProcessInput called with key=%s", keyDesc))
}

// Render creates and returns the current frame with logged messages
func (d *DummyGame) Render() string {
	if d.width <= 0 || d.height <= 0 {
		return "DummyGame: No dimensions set"
	}

	var frame strings.Builder

	// Title line
	title := "=== DUMMY GAME - Method Call Logger ==="
	titlePadding := (d.width - len(title)) / 2
	if titlePadding < 0 {
		titlePadding = 0
	}

	frame.WriteString(strings.Repeat(" ", titlePadding))
	frame.WriteString(title)
	frame.WriteString("\n")

	// Instructions
	instructions := "Press 'q' to quit, any other key to test ProcessInput"
	instrPadding := (d.width - len(instructions)) / 2
	if instrPadding < 0 {
		instrPadding = 0
	}

	frame.WriteString(strings.Repeat(" ", instrPadding))
	frame.WriteString(instructions)
	frame.WriteString("\n")

	// Separator
	frame.WriteString(strings.Repeat("-", d.width))
	frame.WriteString("\n")

	// Display messages (recent first)
	messagesDisplayed := 0
	maxDisplayableMessages := d.height - 5 // Reserve space for title, instructions, separator, and status
	if maxDisplayableMessages < 1 {
		maxDisplayableMessages = 1
	}

	// Show messages in reverse order (most recent first)
	for i := len(d.messages) - 1; i >= 0 && messagesDisplayed < maxDisplayableMessages; i-- {
		message := d.messages[i]
		// Truncate message if it's too long for the width
		if len(message) > d.width {
			message = message[:d.width-3] + "..."
		}
		frame.WriteString(message)
		frame.WriteString("\n")
		messagesDisplayed++
	}

	// Fill remaining space
	remainingLines := d.height - 4 - messagesDisplayed // 4 for title, instructions, separator, status
	for i := 0; i < remainingLines-1; i++ {
		frame.WriteString("\n")
	}

	// Status line at bottom
	status := fmt.Sprintf("Dimensions: %dx%d | Messages: %d/%d",
		d.width, d.height, len(d.messages), d.maxMessages)
	frame.WriteString(status)

	return frame.String()
}

// ShouldQuit returns whether the game should exit
func (d *DummyGame) ShouldQuit() bool {
	return d.shouldQuit
}
