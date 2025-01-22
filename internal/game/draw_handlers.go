package game

import (
	"fmt"
	"image/color"
	"log"

	"github.com/VxVxN/game/internal/shadow"
	"github.com/VxVxN/gamedevlib/raycasting"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (game *Game) drawGameStage(screen *ebiten.Image) {
	if game.sunDirection == shadow.NotSun {
		game.nightImage.Fill(color.Black)
		game.calculateObjects()
		rays := raycasting.RayCasting(game.player.X, game.player.Y, game.objects)

		// Subtract ray triangles from shadow
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Address = ebiten.AddressRepeat
		opt.Blend = ebiten.BlendDestinationOut
		for i, line := range rays {
			nextLine := rays[(i+1)%len(rays)]

			// Draw triangle of area between rays
			v := raycasting.RayVertices(game.player.X, game.player.Y, nextLine.X2, nextLine.Y2, line.X2, line.Y2)
			game.nightImage.DrawTriangles(v, []uint16{0, 1, 2}, game.triangleImage, opt)
		}
	}

	game.background.Draw(screen)
	game.player.Draw(screen)
	game.cars.Draw(screen)
	textFace := &text.GoTextFace{
		Source: game.textFaceSource,
		Size:   24,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(game.windowWidth/2, 0)
	op.ColorScale.Scale(0, 0, 0, 1)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	text.Draw(screen, fmt.Sprintf("Points: %d", int(game.player.Points())), textFace, op)

	if game.sunDirection == shadow.NotSun {
		imageOp := &ebiten.DrawImageOptions{}
		imageOp.ColorScale.ScaleAlpha(0.95)
		screen.DrawImage(game.nightImage, imageOp)
	}
	game.explosionAnimation.Draw(screen)
	if game.stage == GameOverStage {
		textFace = &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   64,
		}

		op = &text.DrawOptions{}
		op.GeoM.Translate(game.windowWidth/2, game.windowHeight/2)
		op.ColorScale.Scale(255, 0, 0, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, "Game over", textFace, op)
	}
}

func (game *Game) drawSetPlayerRecordStage(screen *ebiten.Image) {
	textFace := &text.GoTextFace{
		Source: game.textFaceSource,
		Size:   24,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(game.windowWidth/2, 50)
	op.ColorScale.Scale(255, 255, 255, 1)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	text.Draw(screen, fmt.Sprintf("Your new record: %d", int(game.player.Points())), textFace, op)

	textFace = &text.GoTextFace{
		Source: game.textFaceSource,
		Size:   24,
	}

	op = &text.DrawOptions{}
	op.GeoM.Translate(game.windowWidth/2, 100)
	op.ColorScale.Scale(255, 255, 255, 1)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Enter your name:", textFace, op)
	game.textField.Draw(screen)
}

func (game *Game) drawStatisticsStage(screen *ebiten.Image) {
	records, err := game.statisticer.Load()
	if err != nil {
		log.Fatalf("Failed to load statistics: %v", err)
	}
	textFace := &text.GoTextFace{
		Source: game.textFaceSource,
		Size:   48,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(game.windowWidth/2, 100)
	op.ColorScale.Scale(255, 255, 255, 1)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Player ratings:", textFace, op)

	textFace = &text.GoTextFace{
		Source: game.textFaceSource,
		Size:   24,
	}

	op = &text.DrawOptions{}
	op.GeoM.Translate(game.windowWidth/2, 200)
	op.ColorScale.Scale(255, 255, 255, 1)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Name: Points", textFace, op)

	for i, record := range records {
		textFace = &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   24,
		}

		op = &text.DrawOptions{}
		op.GeoM.Translate(game.windowWidth/2, 250+float64(i*48))
		op.ColorScale.Scale(255, 255, 255, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, fmt.Sprintf("%d) %s: %d", i+1, record.Name, record.Points), textFace, op)
	}
}
