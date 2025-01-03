package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"time"
)

type Game struct {
	width, height int
	y             float64
	road          *ebiten.Image
	globalTime    time.Time
	scrollSpeed   float64
}

func NewGame(width, height int) (*Game, error) {
	road, _, err := ebitenutil.NewImageFromFile("assets/road.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init road image: %v", err)
	}

	ebiten.SetWindowSize(width, height)

	return &Game{
		scrollSpeed: 15.0,
		width:       width,
		height:      height,
		road:        road,
		globalTime:  time.Now(),
	}, nil
}

func (game *Game) Update() error {
	if time.Since(game.globalTime) < time.Second/time.Duration(64) {
		return nil
	}
	game.globalTime = time.Now()
	game.y += game.scrollSpeed
	if game.y > float64(game.road.Bounds().Dy()) {
		game.y = 0
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	for i := -1; i < 5; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(game.width)/2-float64(game.road.Bounds().Dx())/2, float64(game.road.Bounds().Dy()*i)+game.y)
		screen.DrawImage(game.road, op)
	}
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return 1680, 1050
}
