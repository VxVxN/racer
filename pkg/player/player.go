package player

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/VxVxN/game/internal/shadow"
	"github.com/VxVxN/game/pkg/rectangle"
)

type Player struct {
	name   string
	points float64
	*rectangle.Rectangle
	speed  float64
	image  *ebiten.Image
	shadow *shadow.Shadow
	dead   bool
}

func NewPlayer(image *ebiten.Image, shadow *shadow.Shadow, speed float64) *Player {
	return &Player{
		speed:     speed,
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
	player.dead = false
}

func (player *Player) SetName(name string) {
	player.name = name
}

func (player *Player) Name() string {
	return player.name
}

func (player *Player) SetDead(dead bool) {
	player.dead = dead
}

func (player *Player) Dead() bool {
	return player.dead
}

func (player *Player) SetSunDirection(sunDirection shadow.DirectionShadow) {
	player.shadow.SetDirection(sunDirection)
}

func (player *Player) SetSpeed(speed float64) {
	player.speed = speed
}
