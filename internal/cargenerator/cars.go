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
	freeLane     [5]int
}

func New(carImages, truckImages, longTruckImages []*ebiten.Image, screenHeight, startRoad float64, carShadow, truckShadow, longTruckShadow *shadow.Shadow) *CarGenerator {
	carGenerator := &CarGenerator{
		screenHeight: screenHeight,
	}

	carGenerator.cars = make([]*Car, 0, len(carImages))
	for _, image := range carImages {
		car := newCar(image, screenHeight, startRoad, carShadow)
		carGenerator.cars = append(carGenerator.cars, car)
	}
	for _, image := range truckImages {
		car := newCar(image, screenHeight, startRoad, truckShadow)
		carGenerator.cars = append(carGenerator.cars, car)
	}
	for _, image := range longTruckImages {
		car := newCar(image, screenHeight, startRoad, longTruckShadow)
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
	if car.lane != NoLane {
		generator.freeLane[car.lane]--
	}
	for {
		car.Y = float64(-200 - rand.IntN(1800))

		lane := roadLane(rand.IntN(5))
		car.X = car.startRoad + float64(lane)*200 + 65 // 200 - this is the interval between the bands

		if generator.freeLane[lane] == 3 {
			continue
		}
		var availableLane bool
		for i, lineCarCounter := range generator.freeLane {
			if lineCarCounter == 0 && roadLane(i) != lane {
				availableLane = true
				break
			}
		}
		if !availableLane {
			continue
		}
		car.lane = lane
		generator.freeLane[lane]++
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
		if isCollision {
			generator.freeLane[car.lane]--
			continue
		}
		break
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
