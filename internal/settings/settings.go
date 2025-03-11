package settings

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Settings struct {
	SavedSettings *settingValues
	RawSettings   *settingValues
	logger        *slog.Logger
}

type settingValues struct {
	Resolution     Resolution
	MusicVolume    int
	EffectsVolume  int
	CarSensitivity float64
}

type Resolution string

const (
	ResolutionFullScreen Resolution = "Full screen"
	Resolution1920x1080  Resolution = "1920x1080"
	Resolution1680x1050  Resolution = "1680x1050"
	Resolution1280x1024  Resolution = "1280x1024"
	Resolution1280x720   Resolution = "1280x720"
)

func (resolution *Resolution) Size() (int, int) {
	switch *resolution {
	case ResolutionFullScreen:
		return ebiten.Monitor().Size()
	case Resolution1920x1080:
		return 1920, 1080
	case Resolution1680x1050:
		return 1680, 1080
	case Resolution1280x1024:
		return 1024, 1024
	case Resolution1280x720:
		return 1280, 720
	}
	return 1280, 720
}

func New(logger *slog.Logger) (*Settings, error) {
	settings := settingValues{
		Resolution:     ResolutionFullScreen,
		MusicVolume:    100,
		EffectsVolume:  100,
		CarSensitivity: 10,
	}
	data, err := os.ReadFile("settings.json")
	if err == nil {
		if err = json.Unmarshal(data, &settings); err != nil {
			return nil, err
		}
	}
	logger.Info("Started settings", "data", string(data))
	savedSettings := settings
	rawSettings := settings
	return &Settings{
		SavedSettings: &savedSettings,
		RawSettings:   &rawSettings,
		logger:        logger,
	}, nil
}

func (settings *Settings) WriteToFile() error {
	data, err := json.Marshal(settings.SavedSettings)
	if err != nil {
		return err
	}
	if err = os.WriteFile("settings.json", data, 0644); err != nil {
		return err
	}
	settings.logger.Info("Saved settings", "data", string(data))
	return nil
}

func (settings *Settings) Save() {
	savedSettings := *settings.RawSettings
	rawSettings := *settings.RawSettings
	settings.SavedSettings = &savedSettings
	settings.RawSettings = &rawSettings
}
