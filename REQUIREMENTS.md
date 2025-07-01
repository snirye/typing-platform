# ASCII Typing Platform Game - Development Specification

## Project Overview
Create a terminal-based typing game in Go with ASCII graphics where players type words to help a character jump between scrolling platforms.

## Architecture Requirements

### Core Package (`core/`)
- **Purpose**: game engine that handles all game logic and rendering
- **Interface**: 
type GameInterface interface {
	Start(width, height int)
	UpdateDimensions(width, height int)
	ProcessInput(key rune) // No return value - game manages its own state
	Render() string      // Updates game logic and returns rendered frame
	ShouldQuit() bool // indicates exit was preformed
}
- **Responsibilities**: Game state management, collision detection, word validation, scoring, ASCII rendering

### Client Package (`client/`)
- **Purpose**: Platform-specific I/O handling and display
- **Responsibilities**: Terminal setup, keyboard input capture, screen clearing/drawing, game loop
- **Should work**: Terminal
- use https://pkg.go.dev/github.com/nsf/termbox-go for the client

## Game Mechanics

### Core Gameplay Loop
1. **Word Display**: Show target word(s) below current platform
2. **Typing Validation**: Real-time character matching with visual feedback
3. **Success Action**: Character jumps to next platform above on word completion
4. **Platform Scrolling**: Continuous downward movement at configurable speed
5. **Game Over**: When character's platform scrolls below screen bottom
6. **pause, restart, exit**: use `esc` button to pause the game. when paused, use `q` to exit

### Scoring System
- **Score**: Points per completed word (bonus for speed/accuracy)
- **WPM**: Words per minute calculation
- **CPM**: Characters per minute calculation
- **High Score**: Persistent storage between sessions

### Game States
- **Menu**: Start game, view high scores, quit
- **Playing**: Active gameplay
- **Game Over**: Show final stats, play again option
- **Paused**: Optional pause functionality

## Technical Specifications

### ASCII Rendering Requirements
- **Character Sprite**: Simple ASCII character (e.g., `@`, `O`)
- **Platforms**: Horizontal lines using `-` or `=`
- **Current Word**: Highlighted/colored text
- **Typed Characters**: Visual distinction (different color/style)
- **HUD**: Score, WPM, CPM display in corner/border
- **Smooth Animation**: Frame-based movement for scrolling

### Input Handling
- **Alphanumeric**: Word typing
- **Backspace**: Character deletion
- **Escape**: Pause/menu
- **Enter**: Confirm actions in menus
- **Cross-platform**: Work on Windows, macOS, Linux terminals

### Configuration
- **Difficulty Levels**: Adjustable scrolling speed, word complexity
- **Word Lists**: Configurable word sources (common words, programming terms, etc.)
- **Screen Size**: Auto-detect or configurable terminal dimensions

## Implementation Guidelines

### Go-Specific Requirements
- **Modules**: Proper Go module structure with `go.mod`
- **Packages**: Clear separation between `core` and `client`
- **Error Handling**: Proper error propagation and handling
- **Testing**: Unit tests for core game logic
- **Performance**: 60fps target, efficient string operations for ASCII rendering

### Code Organization
```
ascii-type/
├── cmd/
│   └── game/
│       └── main.go          // Entry point
├── internal/
│   ├── core/
│   │   ├── game.go          // Main game logic
│   │   ├── renderer.go      // ASCII rendering
│   │   ├── types.go         // Game structs/interfaces
│   │   └── words.go         // Word management
│   └── client/
│       ├── terminal.go      // Terminal I/O
│       └── display.go       // Screen management
├── assets/
│   └── words.txt           // Word lists
└── go.mod
```

### Development Priorities
1. **Core game engine** with basic rendering
2. **Terminal client** with input handling  
3. **Game mechanics** (jumping, scrolling, word validation)
4. **Scoring system** and statistics
5. **Polish** (menus, animations, difficulty levels)

## Success Criteria
- Runs smoothly in standard terminals (80x24 minimum)
- Responsive input handling with no lag
- Accurate WPM/CPM calculations
- Intuitive gameplay that's immediately understandable
- Clean, maintainable Go code following best practices
