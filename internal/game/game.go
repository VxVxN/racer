package game

import (
	"bytes"
	"fmt"
	"github.com/VxVxN/game/internal/cargenerator"
	"github.com/VxVxN/game/pkg/audioplayer"
	"github.com/VxVxN/game/pkg/background"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/menu"
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
	width, height  float64
	playerCar      *ebiten.Image
	globalTime     time.Time
	scrollSpeed    float64
	textFaceSource *text.GoTextFaceSource
	eventManager   *eventmanager.EventManager
	player         *player.Player
	background     *background.Background
	cars           *cargenerator.CarGenerator
	stage          Stage
	mainMenu       *menu.Menu
	menu           *menu.Menu
}

type Stage int

const (
	GameStage Stage = iota
	GameOverStage
	MainMenuStage
	MenuStage
)

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
	cars := cargenerator.New([]*ebiten.Image{greenCar, orangeCar, redCar, grayCar}, height)

	ebiten.SetWindowSize(int(width), int(height))

	audioPlayer, err := audioplayer.NewAudioPlayer("music")
	if err != nil {
		return nil, fmt.Errorf("failed to init audio player: %v", err)
	}
	audioPlayer.Play()

	textFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, fmt.Errorf("failed to create new face source: %v", err)
	}

	menuTextFace := &text.GoTextFace{
		Source: textFaceSource,
		Size:   32,
	}

	game := &Game{
		scrollSpeed:    20.0,
		width:          width,
		height:         height,
		background:     background.New(road, width),
		playerCar:      playerCar,
		globalTime:     time.Now(),
		eventManager:   eventmanager.NewEventManager(),
		player:         player,
		cars:           cars,
		textFaceSource: textFaceSource,
		stage:          MainMenuStage,
	}

	mainMenu, err := menu.NewMenu(width, height, menuTextFace, []menu.ButtonOptions{
		{
			Text: "New game",
			Action: func() {
				game.stage = GameStage
			},
		},
		{
			Text: "Exit",
			Action: func() {
				os.Exit(0)
			},
		}})
	if err != nil {
		return nil, fmt.Errorf("failed new main menu: %v", err)
	}

	menu, err := menu.NewMenu(width, height, menuTextFace, []menu.ButtonOptions{
		{
			Text: "Continue game",
			Action: func() {
				game.stage = GameStage
			},
		},
		{
			Text: "Go back to the main menu",
			Action: func() {
				game.stage = MainMenuStage
			},
		},
		{
			Text: "Exit",
			Action: func() {
				os.Exit(0)
			},
		}})
	if err != nil {
		return nil, fmt.Errorf("failed new menu: %v", err)
	}

	game.mainMenu = mainMenu
	game.menu = menu
	game.addEvents()

	return game, nil
}

func (game *Game) Update() error {
	game.eventManager.Update()
	if time.Since(game.globalTime) < time.Second/time.Duration(64) {
		return nil
	}
	if game.stage != GameStage {
		return nil
	}
	if game.cars.Collision(game.player.Rectangle) {
		game.stage = GameOverStage
		return nil
	}
	game.globalTime = time.Now()
	game.background.Update(game.scrollSpeed)
	game.cars.Update(game.scrollSpeed - 3)
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	switch game.stage {
	case MainMenuStage:
		game.mainMenu.Draw(screen)
		return
	case MenuStage:
		game.menu.Draw(screen)
		return
	default:
	}
	game.background.Draw(screen)
	game.player.Draw(screen)
	game.cars.Draw(screen)
	if game.stage == GameOverStage {
		textFace := &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   64,
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(game.width/2, game.height/2)
		op.ColorScale.Scale(255, 0, 0, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, "Game over", textFace, op)
	}
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return 1680, 1050
}

func (game *Game) addEvents() {
	game.eventManager.AddPressEvent(ebiten.KeyRight, func() {
		switch game.stage {
		case GameStage:
			if game.player.X < game.width/2+370 {
				game.player.Move(ebiten.KeyRight)
			}
		case GameOverStage:
		case MainMenuStage:
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyLeft, func() {
		switch game.stage {
		case GameStage:
			if game.player.X > game.width/2-480 {
				game.player.Move(ebiten.KeyLeft)
			}
		case GameOverStage:
		case MainMenuStage:
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyUp, func() {
		switch game.stage {
		case GameStage:
			if game.player.Y > 0 {
				game.player.Move(ebiten.KeyUp)
			}
		case GameOverStage:
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyUp, func() {
		switch game.stage {
		case MainMenuStage:
			game.mainMenu.BeforeMenuItem()
		case MenuStage:
			game.menu.BeforeMenuItem()
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyDown, func() {
		switch game.stage {
		case GameStage:
			if game.player.Y < game.height-250 {
				game.player.Move(ebiten.KeyDown)
			}
		case GameOverStage:
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyDown, func() {
		switch game.stage {
		case MainMenuStage:
			game.mainMenu.NextMenuItem()
		case MenuStage:
			game.menu.NextMenuItem()
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyEscape, func() {
		switch game.stage {
		case GameStage:
			game.stage = MenuStage
		case MainMenuStage:
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyEnter, func() {
		switch game.stage {
		case GameStage:
		case GameOverStage:
			game.cars.Reset()
			game.stage = GameStage
		case MainMenuStage:
			game.mainMenu.ClickActiveButton()
		case MenuStage:
			game.menu.ClickActiveButton()
		}
	})
}
