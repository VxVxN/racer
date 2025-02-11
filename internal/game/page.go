package game

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/VxVxN/game/internal/settings"
	"github.com/VxVxN/game/internal/ui"
	"github.com/ebitenui/ebitenui/widget"
)

func newMainPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
	container := ui.NewPageContentContainer()

	buttonOpts := widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
		Position: widget.RowLayoutPositionCenter,
		MaxWidth: 300,
		Stretch:  true,
	}))

	newGameButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("New game", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.TextPadding(res.Button.Padding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.Reset()
		}))
	container.AddChild(newGameButton)

	playerRatingsButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Player ratings", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.setStage(StatisticsStage)
		}))
	container.AddChild(playerRatingsButton)

	settingsButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Settings", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.settingsUI = createUI("Settings", res, newSettingsPage(game, res), false)
			game.setStage(SettingsStage)
		}))
	container.AddChild(settingsButton)

	exitButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Exit", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}))
	container.AddChild(exitButton)

	game.mainMenuButtons = ui.NewButtonControl([]*widget.Button{newGameButton, playerRatingsButton, settingsButton, exitButton})

	return container
}

func newMenuPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
	container := ui.NewPageContentContainer()

	buttonOpts := widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
		Position: widget.RowLayoutPositionCenter,
		MaxWidth: 450,
		Stretch:  true,
	}))

	continueGameButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Continue game", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.TextPadding(res.Button.Padding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.setStage(GameStage)
		}))
	container.AddChild(continueGameButton)

	backToMainMenuButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Go back to the main menu", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.setStage(MainMenuStage)
		}))
	container.AddChild(backToMainMenuButton)

	exitButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Exit", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}))
	container.AddChild(exitButton)

	game.menuButtons = ui.NewButtonControl([]*widget.Button{continueGameButton, backToMainMenuButton, exitButton})

	return container
}

func newPlayerRatingsPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
	container := ui.NewPageContentContainer()

	records, err := game.statisticer.Load()
	if err != nil {
		log.Fatalf("Failed to load statistics: %v", err)
	}

	gridLayoutContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(10, 10))))
	container.AddChild(gridLayoutContainer)

	gridLayoutContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Name", res.Text.TitleFace, res.Text.IdleColor)))

	gridLayoutContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Points", res.Text.TitleFace, res.Text.IdleColor)))

	for _, record := range records {
		gridLayoutContainer.AddChild(widget.NewText(
			widget.TextOpts.Text(record.Name, res.Text.Face, res.Text.IdleColor)))

		gridLayoutContainer.AddChild(widget.NewText(
			widget.TextOpts.Text(strconv.Itoa(record.Points), res.Text.Face, res.Text.IdleColor)))
	}

	return container
}

type setPlayerRatingPage struct {
	widget    widget.PreferredSizeLocateableWidget
	textInput *widget.TextInput
}

func newSetPlayerRatingPage(game *Game, res *ui.UiResources) *setPlayerRatingPage {
	container := ui.NewPageContentContainer()

	gridLayoutContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(10, 10))))
	container.AddChild(gridLayoutContainer)

	gridLayoutContainer.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Your new record: %d", int(game.player.Points())), res.Text.TitleFace, res.Text.IdleColor)))

	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(res.TextInput.Image),
		widget.TextInputOpts.Color(res.TextInput.Color),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(res.TextInput.Face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.TextInput.Face, 2),
		),
	}

	textInput := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.Placeholder("Enter text here"),
		widget.TextInputOpts.AllowDuplicateSubmit(true))...,
	)
	textInput.Focus(true)
	gridLayoutContainer.AddChild(textInput)

	return &setPlayerRatingPage{
		widget:    container,
		textInput: textInput,
	}
}

func newSettingsPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
	container := ui.NewPageContentContainer()

	rayLayoutContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(5))))
	container.AddChild(rayLayoutContainer)

	saveButton := widget.NewButton(
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Save", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.stage = MainMenuStage
			game.settings.Save()
			if err := game.settings.WriteToFile(); err != nil {
				game.logger.Printf("[ERROR] Error saving settings: %v", err)
				return
			}
			game.ApplySettings()
		}))
	rayLayoutContainer.AddChild(saveButton)

	backButton := widget.NewButton(
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Back", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.stage = MainMenuStage
			game.ApplySettings()
		}))
	rayLayoutContainer.AddChild(backButton)

	container.AddChild(ui.NewSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	gridLayoutContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			//widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(10, 10))))
	container.AddChild(gridLayoutContainer)

	gridLayoutContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Resolution", res.Text.Face, res.Text.IdleColor)))

	entries := []interface{}{
		string(settings.ResolutionFullScreen),
		string(settings.Resolution1920x1080),
		string(settings.Resolution1680x1050),
		string(settings.Resolution1280x1024),
		string(settings.Resolution1280x720),
	}

	listResolutionContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Spacing(5))),
	)

	listResolution := ui.NewListComboButton(
		entries,
		func(e interface{}) string {
			return e.(string)
		},
		func(e interface{}) string {
			return e.(string)
		},
		func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			game.settings.RawSettings.Resolution = settings.Resolution(args.Entry.(string))
		},
		res)
	listResolution.SetSelectedEntry(string(game.settings.SavedSettings.Resolution))
	listResolutionContainer.AddChild(listResolution)
	gridLayoutContainer.AddChild(listResolutionContainer)

	sliderMusicVolume := buildSliderMusicVolume(game, res, gridLayoutContainer)
	gridLayoutContainer.AddChild(sliderMusicVolume)

	sliderEffectsVolume := buildSliderEffectsVolume(game, res, gridLayoutContainer)
	gridLayoutContainer.AddChild(sliderEffectsVolume)

	return container
}

func buildSliderMusicVolume(game *Game, res *ui.UiResources, gridLayoutContainer *widget.Container) *widget.Container {
	gridLayoutContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Music Volume", res.Text.Face, res.Text.IdleColor)))

	sliderContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10))),
		widget.ContainerOpts.AutoDisableChildren(),
	)

	var text *widget.Label

	slider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionStart,
		}), widget.WidgetOpts.MinSize(200, 6)),
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.Images(res.Slider.TrackImage, res.Slider.Handle),
		widget.SliderOpts.FixedHandleSize(res.Slider.HandleSize),
		widget.SliderOpts.TrackOffset(5),
		widget.SliderOpts.PageSizeFunc(func() int {
			return 10
		}),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			text.Label = fmt.Sprintf("%d", args.Current)
			game.settings.RawSettings.MusicVolume = args.Current
			game.audioPlayer.SetVolume(float64(game.settings.RawSettings.MusicVolume) / 100)
		}),
	)
	slider.Current = game.settings.SavedSettings.MusicVolume
	sliderContainer.AddChild(slider)

	text = widget.NewLabel(
		widget.LabelOpts.TextOpts(widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionStart,
		}))),
		widget.LabelOpts.Text(fmt.Sprintf("%d", slider.Current), res.Label.Face, res.Label.Text),
	)
	sliderContainer.AddChild(text)
	return sliderContainer
}

func buildSliderEffectsVolume(game *Game, res *ui.UiResources, gridLayoutContainer *widget.Container) *widget.Container {
	gridLayoutContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Effects Volume", res.Text.Face, res.Text.IdleColor)))

	sliderContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10))),
		widget.ContainerOpts.AutoDisableChildren(),
	)

	var text *widget.Label

	slider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionStart,
		}), widget.WidgetOpts.MinSize(200, 6)),
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.Images(res.Slider.TrackImage, res.Slider.Handle),
		widget.SliderOpts.FixedHandleSize(res.Slider.HandleSize),
		widget.SliderOpts.TrackOffset(5),
		widget.SliderOpts.PageSizeFunc(func() int {
			return 10
		}),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			text.Label = fmt.Sprintf("%d", args.Current)
			game.settings.RawSettings.EffectsVolume = args.Current
			game.explosionAnimation.SetVolume(float64(game.settings.RawSettings.EffectsVolume) / 100)
		}),
	)
	slider.Current = game.settings.SavedSettings.EffectsVolume
	sliderContainer.AddChild(slider)

	text = widget.NewLabel(
		widget.LabelOpts.TextOpts(widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionStart,
		}))),
		widget.LabelOpts.Text(fmt.Sprintf("%d", slider.Current), res.Label.Face, res.Label.Text),
	)
	sliderContainer.AddChild(text)
	return sliderContainer
}
