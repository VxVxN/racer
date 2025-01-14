package audioplayer

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"log"
	"math/rand/v2"
	"os"
	"slices"
)

const (
	sampleRate = 48000
)

type AudioPlayer struct {
	audioContext          *audio.Context
	player                *audio.Player
	allMusicFiles         map[string]bool
	notPlayMusicFileNames []string
	logger                *log.Logger
}

func NewAudioPlayer(dirName string, logger *log.Logger) (*AudioPlayer, error) {
	files, err := os.ReadDir(dirName)
	if err != nil {
		return nil, fmt.Errorf("failed to read music directory: %v", err)
	}

	allMusicFiles := make(map[string]bool)
	notPlayMusicFileNames := make([]string, 0, len(files)-1)
	for _, file := range files {
		allMusicFiles[file.Name()] = true
		notPlayMusicFileNames = append(notPlayMusicFileNames, file.Name())
	}

	audioPlayer := &AudioPlayer{
		audioContext:          audio.NewContext(sampleRate),
		allMusicFiles:         allMusicFiles,
		notPlayMusicFileNames: notPlayMusicFileNames,
		logger:                logger,
	}

	return audioPlayer, nil
}

func (audioPlayer *AudioPlayer) Update() error {
	if !audioPlayer.player.IsPlaying() {
		err := audioPlayer.nextMusic()
		if err != nil {
			return err
		}
	}
	return nil
}

func (audioPlayer *AudioPlayer) nextMusic() error {
	if len(audioPlayer.notPlayMusicFileNames) == 0 {
		audioPlayer.notPlayMusicFileNames = make([]string, 0, len(audioPlayer.allMusicFiles))
		for filename, _ := range audioPlayer.allMusicFiles {
			audioPlayer.notPlayMusicFileNames = append(audioPlayer.notPlayMusicFileNames, filename)
		}
	}

	playIndex := rand.IntN(len(audioPlayer.notPlayMusicFileNames))
	fileName := audioPlayer.notPlayMusicFileNames[playIndex]
	audioPlayer.notPlayMusicFileNames = slices.Delete(audioPlayer.notPlayMusicFileNames, playIndex, playIndex+1)

	musicFile, err := os.Open("music/" + fileName) // todo optimize cache file
	if err != nil {
		return fmt.Errorf("failed to open music: %v", err)
	}
	mp3Stream, err := mp3.DecodeF32(musicFile)
	if err != nil {
		return fmt.Errorf("failed to decode music: %v", err)
	}
	audioPlayer.allMusicFiles[fileName] = true

	if audioPlayer.player != nil {
		if err := audioPlayer.player.Close(); err != nil {
			return fmt.Errorf("failed to close audio player: %v", err)
		}
	}
	audioPlayer.logger.Printf("[INFO] Playing music file %s\n", fileName)
	player, err := audioPlayer.audioContext.NewPlayerF32(mp3Stream)
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}
	audioPlayer.player = player
	audioPlayer.player.Play()
	return nil
}

func (audioPlayer *AudioPlayer) Play() {
	if audioPlayer.player == nil {
		_ = audioPlayer.nextMusic()
	}
	audioPlayer.player.Play()
}

func (audioPlayer *AudioPlayer) Pause() {
	audioPlayer.player.Pause()
}
