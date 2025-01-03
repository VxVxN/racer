package game

import (
	"fmt"
	"github.com/VxVxN/game/pkg/audioplayer"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"time"
)

type Game struct {
	width, height float64
	y             float64
	road          *ebiten.Image
	car           *ebiten.Image
	globalTime    time.Time
	scrollSpeed   float64
	eventManager  *eventmanager.EventManager
	player        *player.Player
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

	car := gameElementsSet.SubImage(image.Rect(0, 440, 120, 650)).(*ebiten.Image)

	player := player.NewPlayer(car)

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
		road:         road,
		car:          car,
		globalTime:   time.Now(),
		eventManager: eventmanager.NewEventManager(),
		player:       player,
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
	game.y += game.scrollSpeed
	if game.y > float64(game.road.Bounds().Dy()) {
		game.y = 0
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	for i := -1; i < 5; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(game.width/2-float64(game.road.Bounds().Dx())/2, float64(game.road.Bounds().Dy()*i)+game.y)
		screen.DrawImage(game.road, op)
	}
	game.player.Draw(screen)
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return 1680, 1050
}

func (game *Game) addEvents() {
	game.eventManager.AddPressEvent(ebiten.KeyRight, func() {
		if game.player.X() < game.width/2+float64(game.road.Bounds().Dx())/2-150 {
			game.player.Move(ebiten.KeyRight)
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyLeft, func() {
		if game.player.X() > game.width/2-float64(game.road.Bounds().Dx())/2+40 {
			game.player.Move(ebiten.KeyLeft)
		}
	})
}
