package eventmanager

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EventManager struct {
	eventsPress   map[ebiten.Key][]func()
	eventsPressed map[ebiten.Key][]func()
	defaultEvent  func()
}

func NewEventManager() *EventManager {
	return &EventManager{
		eventsPress:   make(map[ebiten.Key][]func()),
		eventsPressed: make(map[ebiten.Key][]func()),
	}
}

var supportedKeys = []ebiten.Key{
	ebiten.KeyUp,
	ebiten.KeyDown,
	ebiten.KeyLeft,
	ebiten.KeyRight,
	ebiten.KeyTab,
	ebiten.KeyEscape,
	ebiten.KeySpace,
	ebiten.KeyI,
	ebiten.KeyQ,
}

func (eventManager *EventManager) Update() {
	var keyPress, keyPressed ebiten.Key

	for _, supportedKey := range supportedKeys {
		if ebiten.IsKeyPressed(supportedKey) {
			keyPress = supportedKey
		}
		if inpututil.IsKeyJustPressed(supportedKey) {
			keyPressed = supportedKey
		}
	}
	eventsPress, okPress := eventManager.eventsPress[keyPress]
	eventsPressed, okPressed := eventManager.eventsPressed[keyPressed]

	if !okPress && !okPressed && eventManager.defaultEvent != nil {
		eventManager.defaultEvent()
		return // we don't have events
	}
	for _, event := range eventsPress {
		event()
	}
	for _, event := range eventsPressed {
		event()
	}
}

func (eventManager *EventManager) AddPressEvent(key ebiten.Key, event func()) {
	events, _ := eventManager.eventsPress[key]
	events = append(events, event)
	eventManager.eventsPress[key] = events
}

func (eventManager *EventManager) AddPressedEvent(key ebiten.Key, event func()) {
	events, _ := eventManager.eventsPressed[key]
	events = append(events, event)
	eventManager.eventsPressed[key] = events
}

func (eventManager *EventManager) SetDefaultEvent(event func()) {
	eventManager.defaultEvent = event
}
