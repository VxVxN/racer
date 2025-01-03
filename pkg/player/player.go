package player

import "github.com/hajimehoshi/ebiten/v2"

type Player struct {
	x, y  float64
	speed float64
	image *ebiten.Image
}

func NewPlayer(image *ebiten.Image) *Player {
	return &Player{
		speed: 10,
		x:     600,
		y:     950 - float64(image.Bounds().Dy()),
		image: image,
	}
}

func (player *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player.x, player.y)
	screen.DrawImage(player.image, op)
}

func (player *Player) Move(key ebiten.Key) {
	switch key {
	case ebiten.KeyLeft:
		player.x -= player.speed
	case ebiten.KeyRight:
		player.x += player.speed
	case ebiten.KeyUp:
		player.y -= player.speed
	case ebiten.KeyDown:
		player.y += player.speed
	default:
	}
}

func (player *Player) X() float64 {
	return player.x
}

func (player *Player) Y() float64 {
	return player.y
}
