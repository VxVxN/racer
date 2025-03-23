package cargenerator

import (
	"github.com/VxVxN/gamedevlib/rectangle"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/VxVxN/game/internal/shadow"
)

type Car struct {
	screenHeight float64
	startRoad    float64
	*rectangle.Rectangle
	image  *ebiten.Image
	shadow *shadow.Shadow
	lane   roadLane
}

type roadLane int

const (
	NoLane roadLane = iota - 1
	FirstLane
	SecondLane
	ThirdLane
	FourthLane
	FifthLane
)

func newCar(image *ebiten.Image, screenHeight, startRoad float64, shadow *shadow.Shadow) *Car {
	return &Car{
		Rectangle:    rectangle.New(0, 0, float64(image.Bounds().Dx()), float64(image.Bounds().Dy())),
		screenHeight: screenHeight,
		startRoad:    startRoad,
		image:        image,
		shadow:       shadow,
		lane:         NoLane,
	}
}

func (car *Car) Update(scrollSpeed float64) {
	car.Y += scrollSpeed
}

func (car *Car) Draw(screen *ebiten.Image) {
	car.shadow.Draw(screen, car.X, car.Y)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(car.X, car.Y)
	screen.DrawImage(car.image, op)
}

func (car *Car) SetSunDirection(sunDirection shadow.DirectionShadow) {
	car.shadow.SetDirection(sunDirection)
}
