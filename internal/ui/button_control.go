package ui

import (
	"github.com/ebitenui/ebitenui/widget"
)

type ButtonControl struct {
	buttons            []*widget.Button
	currentButtonIndex int
}

func NewButtonControl(buttons []*widget.Button) *ButtonControl {
	buttons[0].Focus(true)
	return &ButtonControl{buttons: buttons}
}

func (bc *ButtonControl) Next() {
	bc.buttons[bc.currentButtonIndex].Focus(false)
	if bc.currentButtonIndex == len(bc.buttons)-1 {
		bc.currentButtonIndex = 0
	} else {
		bc.currentButtonIndex++
	}
	bc.buttons[bc.currentButtonIndex].Focus(true)
}

func (bc *ButtonControl) Before() {
	bc.buttons[bc.currentButtonIndex].Focus(false)
	if bc.currentButtonIndex == 0 {
		bc.currentButtonIndex = len(bc.buttons) - 1
	} else {
		bc.currentButtonIndex--
	}
	bc.buttons[bc.currentButtonIndex].Focus(true)
}

func (bc *ButtonControl) Click() {
	bc.buttons[bc.currentButtonIndex].Click()
}

func (bc *ButtonControl) Pressed() {
	bc.buttons[bc.currentButtonIndex].PressedEvent.Fire(&widget.ButtonPressedEventArgs{Button: bc.buttons[bc.currentButtonIndex]})
}
