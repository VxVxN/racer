package player

import (
	"github.com/VxVxN/game/internal/shadow"
	"github.com/VxVxN/game/pkg/rectangle"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	name   string
	points float64
	*rectangle.Rectangle
	speed  float64
	image  *ebiten.Image
	shadow *shadow.Shadow
}

func NewPlayer(image *ebiten.Image, shadow *shadow.Shadow) *Player {
	return &Player{
		speed:     10,
		Rectangle: rectangle.New(0, 0, float64(image.Bounds().Dx()), float64(image.Bounds().Dy())),
		image:     image,
		shadow:    shadow,
	}
}

func (player *Player) Update() {
	player.points += 0.1
}

func (player *Player) Draw(screen *ebiten.Image) {
	player.shadow.Draw(screen, player.X, player.Y)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player.X, player.Y)
	screen.DrawImage(player.image, op)
}

func (player *Player) SetPosition(x, y float64) {
	player.X = x
	player.Y = y
}

func (player *Player) Move(key ebiten.Key) {
	switch key {
	case ebiten.KeyLeft:
		player.X -= player.speed
	case ebiten.KeyRight:
		player.X += player.speed
	case ebiten.KeyUp:
		player.Y -= player.speed
	case ebiten.KeyDown:
		player.Y += player.speed
	default:
	}
}

func (player *Player) Points() float64 {
	return player.points
}

func (player *Player) Reset() {
	player.points = 0
}

func (player *Player) SetName(name string) {
	player.name = name
}

func (player *Player) Name() string {
	return player.name
}
