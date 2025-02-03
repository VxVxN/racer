package game

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/VxVxN/game/internal/ui"
	"github.com/ebitenui/ebitenui/widget"
)

func mainPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
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

	exitButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text("Exit", res.Button.Face, res.Button.Text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}))
	container.AddChild(exitButton)

	game.mainMenuButtons = ui.NewButtonControl([]*widget.Button{newGameButton, playerRatingsButton, exitButton})

	return container
}

func menuPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
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

func playerRatingsPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
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

func setPlayerRatingPage(game *Game, res *ui.UiResources) widget.PreferredSizeLocateableWidget {
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

	return container
}

func settingsPage(res *ui.UiResources) widget.PreferredSizeLocateableWidget {
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
			// todo save
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
			// todo back
			os.Exit(0)
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
		"Full screen",
		"1920x1080",
		"1680x1050",
		"1280x1024",
		"1280x720",
	}

	cb := ui.NewListComboButton(
		entries,
		func(e interface{}) string {
			return e.(string)
		},
		func(e interface{}) string {
			return e.(string)
		},
		func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			// todo set resolution
		},
		res)
	gridLayoutContainer.AddChild(cb)

	return container
}
