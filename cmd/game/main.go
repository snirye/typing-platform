package main

import (
	"ascii-type/internal/client"
	"ascii-type/internal/core"
	"log"
)

const dummyGame = false // Set to true to use DummyGame for testing

func main() {
	var game core.GameInterface
	var err error
	if dummyGame {
		game = core.NewDummyGame()
	} else {
		game, err = core.NewGame("game_log.txt")
		if err != nil {
			log.Fatalf("Failed to create game: %v", err)
		}
	}

	// Create terminal client
	terminal := client.NewTerminalClient(game)

	// Start the game
	if err := terminal.Run(); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
