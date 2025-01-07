package shadow

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Shadow struct {
	Image           *ebiten.Image
	DirectionShadow DirectionShadow
}
type DirectionShadow int

const (
	SunStraight DirectionShadow = iota
	SunLeft
	SunRight
	SunLeftStraight
	SunRightStraight
)

func New(img *ebiten.Image, direction DirectionShadow) *Shadow {
	return &Shadow{
		Image:           img,
		DirectionShadow: direction,
	}
}

func (shadow *Shadow) Draw(screen *ebiten.Image, x, y float64) {
	var shiftX, shiftY float64
	switch shadow.DirectionShadow {
	case SunStraight:
		shiftY = 40
	case SunLeft:
		shiftX = 40
	case SunRight:
		shiftX = -40
	case SunLeftStraight:
		shiftX = 40
		shiftY = 20
	case SunRightStraight:
		shiftX = -40
		shiftY = 20
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x+shiftX, y+shiftY)
	screen.DrawImage(shadow.Image, op)
}

func (shadow *Shadow) SetDirection(direction DirectionShadow) {
	shadow.DirectionShadow = direction
}
