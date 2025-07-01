# ASCII Typing Platform Game

A terminal-based typing game where you help a character jump between scrolling platforms by typing words correctly.

## Features

- Real-time typing validation with visual feedback
- Scrolling platforms that challenge your speed
- Score tracking with WPM (Words Per Minute) and CPM (Characters Per Minute) statistics
- Pause functionality and game over screen
- Cross-platform terminal support

## Controls

- **Alphanumeric keys**: Type the displayed words
- **Backspace**: Delete the last typed character
- **ESC**: Pause/unpause the game (press 'Q' while paused to quit)
- **Space**: Start game from menu or restart after game over
- **Q**: Quit from menu or when paused
- **Ctrl+C**: Force quit at any time

## How to Play

1. Run the game and press **Space** to start
2. You'll see your character (@) standing on a platform
3. Words appear below platforms - type them correctly to jump to the next platform
4. Platforms continuously scroll down - don't let your character fall off the bottom!
5. Complete words as quickly as possible for higher scores
6. The game ends when your platform scrolls off the bottom of the screen

## Building and Running

```bash
# Build the game
go build ./cmd/game

# Run the game
./game
```

## Requirements

- Go 1.21 or later
- Terminal with at least 80x24 character display
- ANSI color support (most modern terminals)

## Game Mechanics

- **Scoring**: 10 points per character + speed bonus
- **Platform Generation**: New platforms appear as you progress upward
- **Difficulty**: Word length varies to provide appropriate challenge
- **Statistics**: Real-time WPM/CPM calculation and display

## Architecture

The game uses a clean architecture with separation between:
- **Core Engine** (`internal/core/`): Game logic, rendering, and state management
- **Terminal Client** (`internal/client/`): Platform-specific I/O and display handling

This allows for easy porting to different platforms or UI frameworks in the future.

Enjoy improving your typing skills while having fun!
