package main

import (
	"log/slog"
	"os"

	"github.com/nasermirzaei89/gib"
)

func main() {
	err := gib.Run()
	if err != nil {
		slog.Error("Error running game engine", "error", err)
		os.Exit(1)
	}
}
