package main

import (
	"log/slog"
	"os"

	gameengine "github.com/nasermirzaei89/game-engine"
)

func main() {
	err := gameengine.Run()
	if err != nil {
		slog.Error("Error running game engine", "error", err)
		os.Exit(1)
	}
}
