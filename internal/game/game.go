package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"sort"

	"github.com/VxVxN/game/internal/cargenerator"
	"github.com/VxVxN/game/internal/settings"
	"github.com/VxVxN/game/internal/shadow"
	"github.com/VxVxN/game/internal/ui"
	ui2 "github.com/VxVxN/game/internal/ui"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/VxVxN/game/pkg/audioplayer"
	"github.com/VxVxN/game/pkg/background"
	playerpkg "github.com/VxVxN/game/pkg/player"
	"github.com/VxVxN/game/pkg/rectangle"
	"github.com/VxVxN/game/pkg/statisticer"
	"github.com/VxVxN/gamedevlib/eventmanager"
	"github.com/VxVxN/gamedevlib/raycasting"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	// UI
	resourcesUI         *ui.UiResources
	mainMenuUI          *ebitenui.UI
	mainMenuButtons     *ui2.ButtonControl
	menuUI              *ebitenui.UI
	menuButtons         *ui2.ButtonControl
	playerRatingsUI     *ebitenui.UI
	setPlayerRatingUI   *ebitenui.UI
	setPlayerRatingPage *setPlayerRatingPage
	settingsUI          *ebitenui.UI

	windowWidth, windowHeight  float64
	startPlayerX, startPlayerY float64
	scrollSpeed                float64
	textFaceSource             *text.GoTextFaceSource
	eventManager               *eventmanager.EventManager
	player                     *playerpkg.Player
	background                 *background.Background
	cars                       *cargenerator.CarGenerator
	stage                      Stage
	statisticer                *statisticer.Statisticer
	audioPlayer                *audioplayer.AudioPlayer
	nightImage                 *ebiten.Image
	triangleImage              *ebiten.Image
	objects                    []raycasting.Object
	sunDirection               shadow.DirectionShadow
	explosionAnimation         *animation.Animation
	logger                     *log.Logger
	settings                   *settings.Settings
}

type Stage int

const (
	GameStage Stage = iota
	GameOverStage
	MainMenuStage
	MenuStage
	StatisticsStage
	SetPlayerRecordStage
	SettingsStage

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
	case SettingsStage:
		return "SettingsStage"
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

	redTruck := gameElementsSet.SubImage(image.Rect(475, 0, 595, 260)).(*ebiten.Image)
	greenTruck := gameElementsSet.SubImage(image.Rect(600, 0, 720, 260)).(*ebiten.Image)

	blueLongTruck := gameElementsSet.SubImage(image.Rect(760, 0, 900, 425)).(*ebiten.Image)
	greenLongTruck := gameElementsSet.SubImage(image.Rect(900, 0, 1030, 425)).(*ebiten.Image)

	playerShadowImage := vehicleShadowsSet.SubImage(image.Rect(145, 250, 250, 450)).(*ebiten.Image)
	playerShadow := shadow.New(playerShadowImage, shadow.NotSun)

	carShadowImage := vehicleShadowsSet.SubImage(image.Rect(10, 0, 115, 195)).(*ebiten.Image)
	carShadow := shadow.New(carShadowImage, shadow.NotSun)

	truckShadowImage := vehicleShadowsSet.SubImage(image.Rect(140, 0, 255, 245)).(*ebiten.Image)
	truckShadow := shadow.New(truckShadowImage, shadow.NotSun)

	longTruckShadowImage := vehicleShadowsSet.SubImage(image.Rect(280, 0, 385, 445)).(*ebiten.Image)
	longTruckShadow := shadow.New(longTruckShadowImage, shadow.NotSun)

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

	supportedKeys := []ebiten.Key{
		ebiten.KeyUp,
		ebiten.KeyDown,
		ebiten.KeyLeft,
		ebiten.KeyRight,
		ebiten.KeyEscape,
		ebiten.KeyEnter,
	}

	gameSettings, err := settings.New()
	if err != nil {
		return nil, fmt.Errorf("failed to init game settings: %v", err)
	}

	game := &Game{
		scrollSpeed:  10.0,
		windowWidth:  width,
		windowHeight: height,
		background:   background.New(road, width),
		//globalTime:         time.Now(),
		eventManager:       eventmanager.NewEventManager(supportedKeys),
		textFaceSource:     textFaceSource,
		stage:              MainMenuStage,
		statisticer:        statisticer.NewStatisticer(),
		audioPlayer:        audioPlayer,
		nightImage:         ebiten.NewImage(int(width), int(height)),
		triangleImage:      ebiten.NewImage(int(width), int(height)),
		explosionAnimation: animation.NewAnimation(explosionSet, 0, 0, 910, 900, 6),
		cars:               cargenerator.New([]*ebiten.Image{greenCar, orangeCar, redCar, grayCar}, []*ebiten.Image{redTruck, greenTruck}, []*ebiten.Image{blueLongTruck, greenLongTruck}, height, startRoad, carShadow, truckShadow, longTruckShadow),
		player:             playerpkg.NewPlayer(playerCar, playerShadow),
		logger:             logger,
		settings:           gameSettings,
	}

	game.explosionAnimation.SetRepeatable(false)
	game.explosionAnimation.SetScale(0.4, 0.4)
	if err = game.explosionAnimation.SetSound(audioContext, "assets/sounds/silnyiy-vzryiv-starogo-doma.mp3"); err != nil {
		return nil, fmt.Errorf("failed to set sound: %v", err)
	}

	game.ApplySettings()

	m := color.RGBA{ // headlights
		R: 255,
		G: 255,
		B: 255,
		A: 100, // brightness of the headlights
	}
	game.triangleImage.Fill(m)

	res, err := ui.NewUIResources()
	if err != nil {
		return nil, err
	}

	game.resourcesUI = res
	game.mainMenuUI = createUI("Racer", res, newMainPage(game, res), true)
	game.menuUI = createUI("Menu", res, newMenuPage(game, res), true)
	game.playerRatingsUI = createUI("Player ratings", res, newPlayerRatingsPage(game, res), false)
	game.settingsUI = createUI("Settings", res, newSettingsPage(game, res), false)

	game.addEvents()

	return game, nil
}

func (game *Game) Update() error {
	switch game.stage {
	case MainMenuStage:
		game.mainMenuUI.Update()
	case MenuStage:
		game.menuUI.Update()
	case SetPlayerRecordStage:
		game.playerRatingsUI.Update()
	case SettingsStage:
		game.settingsUI.Update()
	}
	game.eventManager.Update()
	//if time.Since(game.globalTime) < time.Second/time.Duration(60) {
	//	return nil
	//}
	if err := game.audioPlayer.Update(); err != nil {
		log.Fatalf("Failed to update audio: %v", err)
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
			game.setPlayerRatingPage = newSetPlayerRatingPage(game, game.resourcesUI)
			game.setPlayerRatingUI = createUI("New record!", game.resourcesUI, game.setPlayerRatingPage.widget, true)
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
	//game.globalTime = time.Now()
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
		game.mainMenuUI.Draw(screen)
	case MenuStage:
		game.menuUI.Draw(screen)
	case StatisticsStage:
		game.playerRatingsUI.Draw(screen)
	case SetPlayerRecordStage:
		game.setPlayerRatingUI.Draw(screen)
	case GameStage, GameOverStage:
		game.drawGameStage(screen)
	case SettingsStage:
		game.settingsUI.Draw(screen)
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
			game.mainMenuButtons.Before()
		case MenuStage:
			game.menuButtons.Before()
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
			game.mainMenuButtons.Next()
		case MenuStage:
			game.menuButtons.Next()
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
			game.mainMenuButtons.Click()
		case MenuStage:
			game.menuButtons.Click()
		case StatisticsStage:
			game.setStage(MainMenuStage)
		case SetPlayerRecordStage:
			if game.setPlayerRatingPage.textInput.GetText() == "" {
				break
			}
			game.player.SetName(game.setPlayerRatingPage.textInput.GetText())
			game.setPlayerRatingPage.textInput.SetText("")

			records, err := game.statisticer.Load()
			if err != nil {
				log.Fatalf("Failed to load statistics: %v", err)
			}
			resultRecords, _ := preparePlayerRatings(records, game.player.Name(), int(game.player.Points()))
			if err := game.statisticer.Save(resultRecords); err != nil {
				log.Fatalf("Failed to save results: %v", err)
			}
			game.playerRatingsUI = createUI("Player ratings", game.resourcesUI, newPlayerRatingsPage(game, game.resourcesUI), false)
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
	sunDirection := shadow.DirectionShadow(1)
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

func createUI(title string, res *ui.UiResources, page widget.PreferredSizeLocateableWidget, center bool) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.TrackHover(false)),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{center, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(res.Background))

	rootContainer.AddChild(ui.HeaderContainer(title, res))

	rootContainer.AddChild(page)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}

func (game *Game) ApplySettings() {
	w, h := game.settings.SavedSettings.Resolution.Size()
	ebiten.SetFullscreen(game.settings.SavedSettings.Resolution == settings.ResolutionFullScreen)
	game.windowWidth, game.windowHeight = float64(w), float64(h)

	if game.settings.SavedSettings.Resolution != settings.ResolutionFullScreen {
		ebiten.SetWindowSize(w, h)
	}

	game.audioPlayer.SetVolume(float64(game.settings.SavedSettings.MusicVolume) / 100)
	game.explosionAnimation.SetVolume(float64(game.settings.SavedSettings.EffectsVolume) / 100)
}
