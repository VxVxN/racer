package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

type Menu struct {
	activeItemMenu            int
	windowWidth, windowHeight float64
	face                      *text.GoTextFace
	buttonOptions             []ButtonOptions
}

type ButtonOptions struct {
	Text   string
	Action func()
}

func NewMenu(windowWidth, windowHeight float64, face *text.GoTextFace, buttonOptions []ButtonOptions) (*Menu, error) {
	return &Menu{
		windowWidth:   windowWidth,
		windowHeight:  windowHeight,
		face:          face,
		buttonOptions: buttonOptions}, nil
}

func (menu *Menu) NextMenuItem() {
	menu.activeItemMenu++
	if menu.activeItemMenu > len(menu.buttonOptions)-1 {
		menu.activeItemMenu = len(menu.buttonOptions) - 1
	}
}

func (menu *Menu) BeforeMenuItem() {
	menu.activeItemMenu--
	if menu.activeItemMenu < 0 {
		menu.activeItemMenu = 0
	}
}

func (menu *Menu) ClickActiveButton() {
	menu.buttonOptions[menu.activeItemMenu].Action()
}

func (menu *Menu) Draw(screen *ebiten.Image) {
	deactivatedButtonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
	activatedButtonColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	for i, buttonOp := range menu.buttonOptions {
		buttonColor := deactivatedButtonColor
		if menu.activeItemMenu == i {
			buttonColor = activatedButtonColor
		}
		op := &text.DrawOptions{}
		op.GeoM.Translate(menu.windowWidth/2, menu.windowHeight/2+float64(i*64))
		op.ColorScale.ScaleWithColor(buttonColor)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, buttonOp.Text, menu.face, op)
	}
}
