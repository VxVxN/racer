package cargenerator

import (
	"math/rand/v2"

	"github.com/VxVxN/game/internal/shadow"
	"github.com/VxVxN/game/pkg/rectangle"
	"github.com/hajimehoshi/ebiten/v2"
)

type CarGenerator struct {
	screenHeight float64
	cars         []*Car
}

func New(images []*ebiten.Image, screenHeight, startRoad float64, carShadow *shadow.Shadow) *CarGenerator {
	carGenerator := &CarGenerator{
		screenHeight: screenHeight,
	}

	carGenerator.cars = make([]*Car, 0, len(images))
	for _, image := range images {
		car := newCar(image, screenHeight, startRoad, carShadow)
		carGenerator.cars = append(carGenerator.cars, car)
	}
	return carGenerator
}

func (generator *CarGenerator) Update(scrollSpeed float64) {
	for i, car := range generator.cars {
		car.Update(scrollSpeed)
		if car.Y > car.screenHeight {
			generator.spawnCar(car, i)
		}
	}
}

func (generator *CarGenerator) spawnCar(car *Car, i int) {
	for {
		car.Y = float64(rand.IntN(400) - 930)
		switch roadLane(rand.IntN(5)) {
		case FirstLane:
			car.X = car.startRoad + 65
		case SecondLane:
			car.X = car.startRoad + 265
		case ThirdLane:
			car.X = car.startRoad + 465
		case FourthLane:
			car.X = car.startRoad + 655
		case FifthLane:
			car.X = car.startRoad + 855
		}
		var isCollision bool
		for j, c := range generator.cars {
			if i == j {
				continue
			}
			if car.Collision(c.Rectangle) {
				isCollision = true
				break
			}
		}
		if !isCollision {
			break
		}
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
	for i, car := range generator.cars {
		generator.spawnCar(car, i)
	}
}

func (generator *CarGenerator) SetSunDirection(sunDirection shadow.DirectionShadow) {
	for _, car := range generator.cars {
		car.SetSunDirection(sunDirection)
	}
}
