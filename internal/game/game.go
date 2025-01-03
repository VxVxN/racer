package game

import (
	"fmt"
	"github.com/VxVxN/game/pkg/audioplayer"
	"github.com/VxVxN/game/pkg/background"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"os"
	"time"
)

type Game struct {
	width, height float64
	playerCar     *ebiten.Image
	car           *ebiten.Image
	globalTime    time.Time
	scrollSpeed   float64
	eventManager  *eventmanager.EventManager
	player        *player.Player
	background    *background.Background
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

	playerCar := gameElementsSet.SubImage(image.Rect(0, 440, 120, 650)).(*ebiten.Image)
	//car := gameElementsSet.SubImage(image.Rect(120, 0, 230, 220)).(*ebiten.Image)

	player := player.NewPlayer(playerCar)

	ebiten.SetWindowSize(int(width), int(height))

	audioPlayer, err := audioplayer.NewAudioPlayer("music")
	if err != nil {
		return nil, fmt.Errorf("failed to init audio player: %v", err)
	}
	audioPlayer.Play()

	game := &Game{
		scrollSpeed:  15.0,
		width:        width,
		height:       height,
		background:   background.New(road, width),
		playerCar:    playerCar,
		globalTime:   time.Now(),
		eventManager: eventmanager.NewEventManager(),
		player:       player,
		//car:          car,
	}

	game.addEvents()

	return game, nil
}

func (game *Game) Update() error {
	game.eventManager.Update()
	if time.Since(game.globalTime) < time.Second/time.Duration(64) {
		return nil
	}
	game.globalTime = time.Now()
	game.background.Update(game.scrollSpeed)
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.background.Draw(screen)
	game.player.Draw(screen)

	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(500, game.y)
	//screen.DrawImage(game.car, op)
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return 1680, 1050
}

func (game *Game) addEvents() {
	game.eventManager.AddPressEvent(ebiten.KeyRight, func() {
		if game.player.X() < game.width/2+370 {
			game.player.Move(ebiten.KeyRight)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyLeft, func() {
		if game.player.X() > game.width/2-480 {
			game.player.Move(ebiten.KeyLeft)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyUp, func() {
		if game.player.Y() > 0 {
			game.player.Move(ebiten.KeyUp)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyDown, func() {
		if game.player.Y() < game.height-250 {
			game.player.Move(ebiten.KeyDown)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyEscape, func() {
		os.Exit(0)
	})
}
