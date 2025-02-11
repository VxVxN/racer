package game

import (
	"fmt"
	"image/color"

	"github.com/VxVxN/game/internal/shadow"
	"github.com/VxVxN/game/internal/stager"
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
	if game.stager.Stage() == stager.GameOverStage {
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
