package game

import (
	"os"

	"github.com/ebitenui/ebitenui/widget"
)

func mainPage(game *Game, res *uiResources) widget.PreferredSizeLocateableWidget {
	container := newPageContentContainer()

	buttonOpts := widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
		Position: widget.RowLayoutPositionCenter,
		MaxWidth: 300,
		Stretch:  true,
	}))

	newGameButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.Text("New game", res.button.face, res.button.text),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.Reset()
		}))
	container.AddChild(newGameButton)

	playerRatingsButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.Text("Player ratings", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			game.setStage(StatisticsStage)
		}))
	container.AddChild(playerRatingsButton)

	exitButton := widget.NewButton(
		buttonOpts,
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.Text("Exit", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}))
	container.AddChild(exitButton)

	game.mainMenuButtons = NewButtonControl([]*widget.Button{newGameButton, playerRatingsButton, exitButton})

	return container
}

func settingsPage(res *uiResources) widget.PreferredSizeLocateableWidget {
	container := newPageContentContainer()

	rayLayoutContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(5))))
	container.AddChild(rayLayoutContainer)

	saveButton := widget.NewButton(
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.Text("Save", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			// todo save
		}))
	rayLayoutContainer.AddChild(saveButton)

	backButton := widget.NewButton(
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.Text("Back", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			// todo back
			os.Exit(0)
		}))
	rayLayoutContainer.AddChild(backButton)

	container.AddChild(newSeparator(res, widget.RowLayoutData{
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
		widget.TextOpts.Text("Resolution", res.text.face, res.text.idleColor)))

	entries := []interface{}{
		"Full screen",
		"1920x1080",
		"1680x1050",
		"1280x1024",
		"1280x720",
	}

	cb := newListComboButton(
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
