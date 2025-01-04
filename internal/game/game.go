package game

import (
	"bytes"
	"fmt"
	cars2 "github.com/VxVxN/game/internal/cargenerator"
	"github.com/VxVxN/game/pkg/audioplayer"
	"github.com/VxVxN/game/pkg/background"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"os"
	"time"
)

type Game struct {
	width, height float64
	playerCar     *ebiten.Image
	globalTime    time.Time
	scrollSpeed   float64
	gameOverFace  *text.GoTextFaceSource
	eventManager  *eventmanager.EventManager
	player        *player.Player
	background    *background.Background
	cars          *cars2.CarGenerator
	gameOver      bool
}

func NewGame(width, height float64) (*Game, error) {
	road, _, err := ebitenutil.NewImageFromFile("assets/road.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init road image: %v", err)
	}

	gameElementsSet, _, err := ebitenutil.NewImageFromFile("assets/game elements.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init game elements image: %v", err)
	}

	playerCar := gameElementsSet.SubImage(image.Rect(0, 450, 110, 650)).(*ebiten.Image)
	greenCar := gameElementsSet.SubImage(image.Rect(0, 0, 110, 220)).(*ebiten.Image)
	orangeCar := gameElementsSet.SubImage(image.Rect(120, 0, 230, 220)).(*ebiten.Image)
	redCar := gameElementsSet.SubImage(image.Rect(240, 0, 350, 220)).(*ebiten.Image)
	grayCar := gameElementsSet.SubImage(image.Rect(360, 0, 470, 220)).(*ebiten.Image)

	player := player.NewPlayer(playerCar)
	cars := cars2.New([]*ebiten.Image{greenCar, orangeCar, redCar, grayCar}, height)

	ebiten.SetWindowSize(int(width), int(height))

	audioPlayer, err := audioplayer.NewAudioPlayer("music")
	if err != nil {
		return nil, fmt.Errorf("failed to init audio player: %v", err)
	}
	audioPlayer.Play()

	gameOverFace, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, fmt.Errorf("failed to create new face for game ofer face(font): %v", err)
	}

	game := &Game{
		scrollSpeed:  20.0,
		width:        width,
		height:       height,
		background:   background.New(road, width),
		playerCar:    playerCar,
		globalTime:   time.Now(),
		eventManager: eventmanager.NewEventManager(),
		player:       player,
		cars:         cars,
		gameOverFace: gameOverFace,
	}

	game.addEvents()

	return game, nil
}

func (game *Game) Update() error {
	game.eventManager.Update()
	if time.Since(game.globalTime) < time.Second/time.Duration(64) {
		return nil
	}
	if game.cars.Collision(game.player.Rectangle) {
		game.gameOver = true
		return nil
	}
	game.globalTime = time.Now()
	game.background.Update(game.scrollSpeed)
	game.cars.Update(game.scrollSpeed - 3)
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.background.Draw(screen)
	game.player.Draw(screen)
	game.cars.Draw(screen)
	if game.gameOver {
		op := &text.DrawOptions{}
		f := &text.GoTextFace{
			Source: game.gameOverFace,
			Size:   64,
		}
		op.GeoM.Translate(game.width/2-150, game.height/2)
		op.ColorScale.Scale(255, 0, 0, 1)
		text.Draw(screen, "Game over", f, op)
	}
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return 1680, 1050
}

func (game *Game) addEvents() {
	game.eventManager.AddPressEvent(ebiten.KeyRight, func() {
		if game.gameOver {
			return
		}
		if game.player.X < game.width/2+370 {
			game.player.Move(ebiten.KeyRight)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyLeft, func() {
		if game.gameOver {
			return
		}
		if game.player.X > game.width/2-480 {
			game.player.Move(ebiten.KeyLeft)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyUp, func() {
		if game.gameOver {
			return
		}
		if game.player.Y > 0 {
			game.player.Move(ebiten.KeyUp)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyDown, func() {
		if game.gameOver {
			return
		}
		if game.player.Y < game.height-250 {
			game.player.Move(ebiten.KeyDown)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyEscape, func() {
		os.Exit(0)
	})
	game.eventManager.AddPressEvent(ebiten.KeyEnter, func() {
		game.cars.Reset()
		game.gameOver = false
	})
}
