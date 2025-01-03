package main

import (
	"github.com/VxVxN/game/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game, err := game.NewGame(1680, 1050)
	if err != nil {
		log.Fatalf("Failed to init game: %v", err)
	}

	ebiten.SetWindowTitle("Racer")

	if err = ebiten.RunGame(game); err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}
