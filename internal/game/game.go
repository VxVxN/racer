package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand/v2"
	"os"
	"sort"
	"time"

	"github.com/VxVxN/game/internal/cargenerator"
	"github.com/VxVxN/game/internal/shadow"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/VxVxN/game/pkg/audioplayer"
	"github.com/VxVxN/game/pkg/background"
	playerpkg "github.com/VxVxN/game/pkg/player"
	"github.com/VxVxN/game/pkg/rectangle"
	"github.com/VxVxN/game/pkg/statisticer"
	"github.com/VxVxN/game/pkg/textfield"
	"github.com/VxVxN/gamedevlib/eventmanager"
	"github.com/VxVxN/gamedevlib/menu"
	"github.com/VxVxN/gamedevlib/raycasting"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	windowWidth, windowHeight  float64
	startPlayerX, startPlayerY float64
	globalTime                 time.Time
	scrollSpeed                float64
	textFaceSource             *text.GoTextFaceSource
	eventManager               *eventmanager.EventManager
	player                     *playerpkg.Player
	background                 *background.Background
	cars                       *cargenerator.CarGenerator
	stage                      Stage
	mainMenu                   *menu.Menu
	menu                       *menu.Menu
	statisticer                *statisticer.Statisticer
	textField                  *textfield.TextField
	audioPlayer                *audioplayer.AudioPlayer
	nightImage                 *ebiten.Image
	triangleImage              *ebiten.Image
	objects                    []raycasting.Object
	sunDirection               shadow.DirectionShadow
	explosionAnimation         *animation.Animation
	logger                     *log.Logger
}

type Stage int

const (
	GameStage Stage = iota
	GameOverStage
	MainMenuStage
	MenuStage
	StatisticsStage
	SetPlayerRecordStage

	sampleRate = 48000
)

func (s Stage) String() string {
	switch s {
	case GameStage:
		return "GameStage"
	case GameOverStage:
		return "GameOverStage"
	case MainMenuStage:
		return "MainMenuStage"
	case MenuStage:
		return "MenuStage"
	case StatisticsStage:
		return "StatisticsStage"
	case SetPlayerRecordStage:
		return "SetPlayerRecordStage"
	}
	return ""
}

func NewGame() (*Game, error) {
	logger := log.Default()
	w, h := ebiten.Monitor().Size()
	logger.Printf("[INFO] Monitor size(%dx%d)", w, h)
	width, height := float64(w), float64(h)

	road, _, err := ebitenutil.NewImageFromFile("assets/road.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init road image: %v", err)
	}

	gameElementsSet, _, err := ebitenutil.NewImageFromFile("assets/game elements.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init game elements image: %v", err)
	}

	vehicleShadowsSet, _, err := ebitenutil.NewImageFromFile("assets/vehicleShadows.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init game vehicle shadows image: %v", err)
	}

	explosionSet, _, err := ebitenutil.NewImageFromFile("assets/explosion.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init game explosion image: %v", err)
	}

	playerCar := gameElementsSet.SubImage(image.Rect(0, 450, 110, 650)).(*ebiten.Image)
	greenCar := gameElementsSet.SubImage(image.Rect(0, 0, 110, 210)).(*ebiten.Image)
	orangeCar := gameElementsSet.SubImage(image.Rect(120, 0, 230, 210)).(*ebiten.Image)
	redCar := gameElementsSet.SubImage(image.Rect(240, 0, 350, 210)).(*ebiten.Image)
	grayCar := gameElementsSet.SubImage(image.Rect(360, 0, 470, 210)).(*ebiten.Image)

	playerShadowImage := vehicleShadowsSet.SubImage(image.Rect(145, 250, 250, 450)).(*ebiten.Image)
	playerShadow := shadow.New(playerShadowImage, shadow.NotSun)

	carShadowImage := vehicleShadowsSet.SubImage(image.Rect(10, 0, 115, 195)).(*ebiten.Image)
	carShadow := shadow.New(carShadowImage, shadow.NotSun)

	startRoad := width/2 - float64(road.Bounds().Dx())/2

	ebiten.SetWindowSize(int(width), int(height))

	audioContext := audio.NewContext(sampleRate)

	audioPlayer, err := audioplayer.NewAudioPlayer(audioContext, "music", logger)
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

	inputTextFace := &text.GoTextFace{
		Source: textFaceSource,
		Size:   32,
	}

	textField := textfield.NewTextField(inputTextFace, image.Rect(16, 150, int(width-16), int(150+inputTextFace.Size)), false)

	supportedKeys := []ebiten.Key{
		ebiten.KeyUp,
		ebiten.KeyDown,
		ebiten.KeyLeft,
		ebiten.KeyRight,
		ebiten.KeyEscape,
		ebiten.KeyEnter,
	}

	game := &Game{
		scrollSpeed:        20.0,
		windowWidth:        width,
		windowHeight:       height,
		background:         background.New(road, width),
		globalTime:         time.Now(),
		eventManager:       eventmanager.NewEventManager(supportedKeys),
		textFaceSource:     textFaceSource,
		stage:              MainMenuStage,
		statisticer:        statisticer.NewStatisticer(),
		textField:          textField,
		audioPlayer:        audioPlayer,
		nightImage:         ebiten.NewImage(int(width), int(height)),
		triangleImage:      ebiten.NewImage(int(width), int(height)),
		explosionAnimation: animation.NewAnimation(explosionSet, 0, 0, 910, 900, 6),
		cars:               cargenerator.New([]*ebiten.Image{greenCar, orangeCar, redCar, grayCar}, height, startRoad, carShadow),
		player:             playerpkg.NewPlayer(playerCar, playerShadow),
		logger:             logger,
	}
	game.explosionAnimation.SetRepeatable(false)
	game.explosionAnimation.SetScale(0.4, 0.4)
	if err = game.explosionAnimation.SetSound(audioContext, "assets/sounds/silnyiy-vzryiv-starogo-doma.mp3"); err != nil {
		return nil, fmt.Errorf("failed to set sound: %v", err)
	}

	m := color.RGBA{ // headlights
		R: 255,
		G: 255,
		B: 255,
		A: 100, // brightness of the headlights
	}
	game.triangleImage.Fill(m)

	mainMenu, err := menu.NewMenu(width, height, menuTextFace, []menu.ButtonOptions{
		{
			Text: "New game",
			Action: func() {
				game.Reset()
			},
		},
		{
			Text: "Player ratings",
			Action: func() {
				game.setStage(StatisticsStage)
			},
		},
		{
			Text: "Exit",
			Action: func() {
				os.Exit(0)
			},
		}}, menu.MenuOptions{
		ButtonPadding: 20,
	})
	if err != nil {
		return nil, fmt.Errorf("failed new main menu: %v", err)
	}

	menu, err := menu.NewMenu(width, height, menuTextFace, []menu.ButtonOptions{
		{
			Text: "Continue game",
			Action: func() {
				game.setStage(GameStage)
			},
		},
		{
			Text: "Go back to the main menu",
			Action: func() {
				game.setStage(MainMenuStage)
			},
		},
		{
			Text: "Exit",
			Action: func() {
				os.Exit(0)
			},
		}}, menu.MenuOptions{
		ButtonPadding: 20,
	})
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
	if err := game.audioPlayer.Update(); err != nil {
		log.Fatalf("Failed to update audio: %v", err)
	}
	if game.stage == SetPlayerRecordStage {
		game.textField.Focus()
		if err := game.textField.Update(); err != nil {
			log.Fatalf("Failed to update text: %v", err)
		}
		return nil
	}
	if game.stage != GameStage {
		return nil
	}

	game.explosionAnimation.Update(0.1)

	if game.player.Dead() {
		return nil
	}

	if game.cars.Collision(game.player.Rectangle) {
		game.audioPlayer.Pause()
		game.player.SetDead(true)
		game.logger.Println("[DBG] Collision detected")
		game.explosionAnimation.SetPosition(game.player.X*2.15, game.player.Y*2.15)
		game.explosionAnimation.Start()
		game.explosionAnimation.SetCallback(func() {
			defer game.audioPlayer.Play()
			game.setStage(GameOverStage)
			records, err := game.statisticer.Load()
			if err != nil {
				log.Fatalf("Failed to load statistics: %v", err)
			}
			_, isRecord := preparePlayerRatings(records, game.player.Name(), int(game.player.Points()))
			if !isRecord {
				return
			}
			game.setStage(SetPlayerRecordStage)
		})
		return nil
	}
	game.globalTime = time.Now()
	game.background.Update(game.scrollSpeed)
	game.player.Update()
	game.cars.Update(game.scrollSpeed - 3)
	return nil
}

func preparePlayerRatings(records []statisticer.Record, playerName string, playerPoints int) ([]statisticer.Record, bool) {
	if len(records) == 0 {
		return append(records, statisticer.NewRecord(playerName, playerPoints)), true
	}

	var isRecord bool
	for _, record := range records {
		if playerPoints > record.Points {
			isRecord = true
			break
		}
	}
	if len(records) < 10 {
		isRecord = true
	}
	if !isRecord {
		return records, false
	}
	records = append(records, statisticer.NewRecord(playerName, playerPoints))

	sort.Slice(records, func(i, j int) bool {
		return records[i].Points > records[j].Points
	})
	if len(records) > 10 {
		records = records[:10]
	}
	return records, true
}

func (game *Game) Draw(screen *ebiten.Image) {
	switch game.stage {
	case MainMenuStage:
		game.mainMenu.Draw(screen)
	case MenuStage:
		game.menu.Draw(screen)
	case StatisticsStage:
		game.drawStatisticsStage(screen)
	case SetPlayerRecordStage:
		game.drawSetPlayerRecordStage(screen)
	case GameStage, GameOverStage:
		game.drawGameStage(screen)
	default:
	}
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return screenWidthPx, screenHeightPx
}

func (game *Game) addEvents() {
	game.eventManager.AddPressEvent(ebiten.KeyRight, func() {
		switch game.stage {
		case GameStage:
			if !game.player.Dead() && game.player.X < game.windowWidth/2+370 {
				game.player.Move(ebiten.KeyRight)
			}
		case GameOverStage:
		case MainMenuStage:
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyLeft, func() {
		switch game.stage {
		case GameStage:
			if !game.player.Dead() && game.player.X > game.windowWidth/2-480 {
				game.player.Move(ebiten.KeyLeft)
			}
		case GameOverStage:
		case MainMenuStage:
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyUp, func() {
		switch game.stage {
		case GameStage:
			if !game.player.Dead() && game.player.Y > 0 {
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
			if !game.player.Dead() && game.player.Y < game.windowHeight-210 {
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
		case GameStage, GameOverStage:
			game.setStage(MenuStage)
		case StatisticsStage:
			game.setStage(MainMenuStage)
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyEnter, func() {
		switch game.stage {
		case GameStage:
		case GameOverStage:
			game.Reset()
		case MainMenuStage:
			game.mainMenu.ClickActiveButton()
		case MenuStage:
			game.menu.ClickActiveButton()
		case StatisticsStage:
			game.setStage(MainMenuStage)
		case SetPlayerRecordStage:
			game.player.SetName(game.textField.Text())
			records, err := game.statisticer.Load()
			if err != nil {
				log.Fatalf("Failed to load statistics: %v", err)
			}
			resultRecords, _ := preparePlayerRatings(records, game.player.Name(), int(game.player.Points()))
			if err := game.statisticer.Save(resultRecords); err != nil {
				log.Fatalf("Failed to save results: %v", err)
			}
			game.setStage(StatisticsStage)
		}
	})
}

func (game *Game) calculateObjects() {
	game.objects = []raycasting.Object{
		ConvertRectangleToObject(*rectangle.New(0, 0, game.windowWidth, game.windowHeight)),
		*raycasting.NewObject([]raycasting.Line{{ // right ray
			float64(game.player.X) + 100,
			float64(game.player.Y) + 2,
			game.player.X - game.startPlayerX + game.windowWidth - 500,
			game.player.Y - game.windowHeight}}),
		*raycasting.NewObject([]raycasting.Line{{ // left ray
			game.player.X,
			float64(game.player.Y) + 2,
			game.player.X - game.startPlayerX + 580,
			game.player.Y - game.windowHeight}}),
		*raycasting.NewObject([]raycasting.Line{{0, game.player.Y, game.windowWidth, float64(game.player.Y) + 2}}),
	}
}

func (game *Game) setStage(newStage Stage) {
	game.stage = newStage
	game.logger.Printf("[INFO] Setting stage to %v", newStage)
}

func ConvertRectangleToObject(rectangle rectangle.Rectangle) raycasting.Object {
	return raycasting.Object{
		Walls: []raycasting.Line{
			{
				X1: rectangle.X,
				Y1: rectangle.Y,
				X2: rectangle.X,
				Y2: rectangle.Y + rectangle.Height,
			},
			{
				X1: rectangle.X,
				Y1: rectangle.Y + rectangle.Height,
				X2: rectangle.X + rectangle.Width,
				Y2: rectangle.Y + rectangle.Height,
			},
			{
				X1: rectangle.X + rectangle.Width,
				Y1: rectangle.Y + rectangle.Height,
				X2: rectangle.X + rectangle.Width,
				Y2: rectangle.Y,
			},
			{
				X1: rectangle.X + rectangle.Width,
				Y1: rectangle.Y,
				X2: rectangle.X,
				Y2: rectangle.Y,
			},
		},
	}
}

func (game *Game) Reset() {
	sunDirection := shadow.DirectionShadow(rand.IntN(6))
	game.sunDirection = sunDirection

	game.setStage(GameStage)
	game.player.Reset()
	game.player.SetPosition(game.windowWidth/2-float64(game.player.Rectangle.Width)/2, game.windowHeight/2)
	game.player.SetSunDirection(sunDirection)
	game.startPlayerX = game.player.X
	game.startPlayerY = game.player.Y

	game.cars.Reset()
	game.cars.SetSunDirection(sunDirection)
	game.explosionAnimation.Reset()
}
