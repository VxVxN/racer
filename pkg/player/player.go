package player

import "github.com/hajimehoshi/ebiten/v2"

type Player struct {
	x     float64
	speed float64
	image *ebiten.Image
}

func NewPlayer(image *ebiten.Image) *Player {
	return &Player{
		speed: 10,
		x:     600,
		image: image,
	}
}

func (player *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player.x, 1050-float64(player.image.Bounds().Dy())-100)
	screen.DrawImage(player.image, op)
}

func (player *Player) Move(key ebiten.Key) {
	switch key {
	case ebiten.KeyLeft:
		player.x -= player.speed
	case ebiten.KeyRight:
		player.x += player.speed
	default:
	}
}

func (player *Player) X() float64 {
	return player.x
}
