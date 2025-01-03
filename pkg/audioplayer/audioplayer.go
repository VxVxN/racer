package audioplayer

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"math/rand/v2"
	"os"
)

const (
	sampleRate = 48000
)

type AudioPlayer struct {
	player *audio.Player
}

func NewAudioPlayer(dirName string) (*AudioPlayer, error) {
	files, err := os.ReadDir(dirName)
	if err != nil {
		return nil, fmt.Errorf("failed to read music directory: %v", err)
	}

	fileName := files[rand.IntN(len(files))].Name()

	musicFile, err := os.Open("music/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open music: %v", err)
	}
	streamMusic, err := mp3.DecodeF32(musicFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode music: %v", err)
	}
	audioContext := audio.NewContext(sampleRate)
	player, err := audioContext.NewPlayerF32(streamMusic)
	if err != nil {
		return nil, err
	}
	return &AudioPlayer{
		player: player,
	}, nil
}

func (audioPlayer *AudioPlayer) Play() {
	audioPlayer.player.Play()
}
