package main

import (
	"log"

	"github.com/VxVxN/game/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game, err := game.NewGame()
	if err != nil {
		log.Fatalf("Failed to init game: %v", err)
	}

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Racer")

	if err = ebiten.RunGame(game); err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}
