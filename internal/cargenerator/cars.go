package cargenerator

import (
	"github.com/VxVxN/game/pkg/rectangle"
	"github.com/hajimehoshi/ebiten/v2"
)

type CarGenerator struct {
	screenHeight float64
	cars         []*car
}

func New(images []*ebiten.Image, screenHeight float64) *CarGenerator {
	cars := make([]*car, 0, len(images))
	for _, image := range images {
		cars = append(cars, newCar(image, screenHeight))
	}
	return &CarGenerator{
		cars:         cars,
		screenHeight: screenHeight,
	}
}

func (generator *CarGenerator) Update(scrollSpeed float64) {
	for _, car := range generator.cars {
		car.Update(scrollSpeed)
	}
}

func (generator *CarGenerator) Draw(screen *ebiten.Image) {
	for _, car := range generator.cars {
		car.Draw(screen)
	}
}

func (generator *CarGenerator) Collision(rectangle *rectangle.Rectangle) bool {
	for _, car := range generator.cars {
		if car.Collision(rectangle) {
			return true
		}
	}
	return false
}

func (generator *CarGenerator) Reset() {
	for _, car := range generator.cars {
		car.Reset()
	}
}
