package main

import (
	"ascii-type/internal/client"
	"ascii-type/internal/core"
	"log"
)

func main() {
	// Create game engine
	game := core.NewGame()

	// Create terminal client
	terminal := client.NewTerminalClient(game)

	// Start the game
	if err := terminal.Run(); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
