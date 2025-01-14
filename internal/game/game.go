package game

import (
	"bytes"
	"fmt"
	"github.com/VxVxN/game/internal/cargenerator"
	"github.com/VxVxN/game/internal/shadow"
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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"log"
	"math/rand/v2"
	"os"
	"sort"
	"time"
)

type Game struct {
	width, height              float64
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
}

type Stage int

const (
	GameStage Stage = iota
	GameOverStage
	MainMenuStage
	MenuStage
	StatisticsStage
	SetPlayerRecordStage
)

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

	audioPlayer, err := audioplayer.NewAudioPlayer("music", logger)
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
		scrollSpeed:    20.0,
		width:          width,
		height:         height,
		background:     background.New(road, width),
		globalTime:     time.Now(),
		eventManager:   eventmanager.NewEventManager(supportedKeys),
		textFaceSource: textFaceSource,
		stage:          MainMenuStage,
		statisticer:    statisticer.NewStatisticer(),
		textField:      textField,
		audioPlayer:    audioPlayer,
		nightImage:     ebiten.NewImage(int(width), int(height)),
		triangleImage:  ebiten.NewImage(int(width), int(height)),
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
				sunDirection := shadow.DirectionShadow(rand.IntN(6))
				playerShadow.SetDirection(sunDirection)
				carShadow.SetDirection(sunDirection)
				game.sunDirection = sunDirection

				game.stage = GameStage
				game.player = playerpkg.NewPlayer(playerCar, playerShadow)
				game.player.SetPosition(width/2-float64(playerCar.Bounds().Dx())/2, height/2)
				game.startPlayerX = game.player.X
				game.startPlayerY = game.player.Y

				game.cars = cargenerator.New([]*ebiten.Image{greenCar, orangeCar, redCar, grayCar}, height, startRoad, carShadow)
			},
		},
		{
			Text: "Player ratings",
			Action: func() {
				game.stage = StatisticsStage
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
	if game.cars.Collision(game.player.Rectangle) {
		game.stage = GameOverStage
		records, err := game.statisticer.Load()
		if err != nil {
			log.Fatalf("Failed to load statistics: %v", err)
		}
		_, isRecord := preparePlayerRatings(records, game.player.Name(), int(game.player.Points()))
		if !isRecord {
			return nil
		}
		game.stage = SetPlayerRecordStage
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
		return
	case MenuStage:
		game.menu.Draw(screen)
		return
	case StatisticsStage:
		records, err := game.statisticer.Load()
		if err != nil {
			log.Fatalf("Failed to load statistics: %v", err)
		}
		textFace := &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   48,
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(game.width/2, 100)
		op.ColorScale.Scale(255, 255, 255, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, "Player ratings:", textFace, op)

		textFace = &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   24,
		}

		op = &text.DrawOptions{}
		op.GeoM.Translate(game.width/2, 200)
		op.ColorScale.Scale(255, 255, 255, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, "Name: Points", textFace, op)

		for i, record := range records {
			textFace = &text.GoTextFace{
				Source: game.textFaceSource,
				Size:   24,
			}

			op = &text.DrawOptions{}
			op.GeoM.Translate(game.width/2, 250+float64(i*48))
			op.ColorScale.Scale(255, 255, 255, 1)
			op.LayoutOptions.PrimaryAlign = text.AlignCenter
			text.Draw(screen, fmt.Sprintf("%d) %s: %d", i+1, record.Name, record.Points), textFace, op)
		}
		return
	case SetPlayerRecordStage:
		textFace := &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   24,
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(game.width/2, 50)
		op.ColorScale.Scale(255, 255, 255, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, fmt.Sprintf("Your new record: %d", int(game.player.Points())), textFace, op)

		textFace = &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   24,
		}

		op = &text.DrawOptions{}
		op.GeoM.Translate(game.width/2, 100)
		op.ColorScale.Scale(255, 255, 255, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, "Enter your name:", textFace, op)
		game.textField.Draw(screen)
	case GameStage, GameOverStage:
		if game.sunDirection == shadow.NotSun {
			game.nightImage.Fill(color.Black)
			game.calculateObjects()
			rays := raycasting.RayCasting(game.player.X, game.player.Y, game.objects)

			// Subtract ray triangles from shadow
			opt := &ebiten.DrawTrianglesOptions{}
			opt.Address = ebiten.AddressRepeat
			opt.Blend = ebiten.BlendDestinationOut
			for i, line := range rays {
				nextLine := rays[(i+1)%len(rays)]

				// Draw triangle of area between rays
				v := raycasting.RayVertices(game.player.X, game.player.Y, nextLine.X2, nextLine.Y2, line.X2, line.Y2)
				game.nightImage.DrawTriangles(v, []uint16{0, 1, 2}, game.triangleImage, opt)
			}
		}

		game.background.Draw(screen)
		game.player.Draw(screen)
		game.cars.Draw(screen)
		textFace := &text.GoTextFace{
			Source: game.textFaceSource,
			Size:   24,
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(game.width/2, 0)
		op.ColorScale.Scale(0, 0, 0, 1)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, fmt.Sprintf("Points: %d", int(game.player.Points())), textFace, op)

		if game.sunDirection == shadow.NotSun {
			imageOp := &ebiten.DrawImageOptions{}
			imageOp.ColorScale.ScaleAlpha(0.95)
			screen.DrawImage(game.nightImage, imageOp)
		}

		if game.stage == GameOverStage {
			textFace = &text.GoTextFace{
				Source: game.textFaceSource,
				Size:   64,
			}

			op = &text.DrawOptions{}
			op.GeoM.Translate(game.width/2, game.height/2)
			op.ColorScale.Scale(255, 0, 0, 1)
			op.LayoutOptions.PrimaryAlign = text.AlignCenter
			text.Draw(screen, "Game over", textFace, op)
		}
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
			if game.player.Y < game.height-210 {
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
			game.stage = MenuStage
		case StatisticsStage:
			game.stage = MainMenuStage
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyEnter, func() {
		switch game.stage {
		case GameStage:
		case GameOverStage:
			game.player.Reset()
			game.cars.Reset()
			game.stage = GameStage
		case MainMenuStage:
			game.mainMenu.ClickActiveButton()
		case MenuStage:
			game.menu.ClickActiveButton()
		case StatisticsStage:
			game.stage = MainMenuStage
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
			game.stage = StatisticsStage
		}
	})
}

func (game *Game) calculateObjects() {
	game.objects = []raycasting.Object{
		ConvertRectangleToObject(*rectangle.New(0, 0, game.width, game.height)),
		*raycasting.NewObject([]raycasting.Line{{ // right ray
			float64(game.player.X) + 100,
			float64(game.player.Y) + 2,
			game.player.X - game.startPlayerX + game.width - 500,
			game.player.Y - game.height}}),
		*raycasting.NewObject([]raycasting.Line{{ // left ray
			game.player.X,
			float64(game.player.Y) + 2,
			game.player.X - game.startPlayerX + 580,
			game.player.Y - game.height}}),
		*raycasting.NewObject([]raycasting.Line{{0, game.player.Y, game.width, float64(game.player.Y) + 2}}),
	}
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
