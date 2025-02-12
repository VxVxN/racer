package audioplayer

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type AudioPlayer struct {
	audioContext     *audio.Context
	player           *audio.Player
	currentSongName  string
	allMusicFiles    []*os.File
	currentSongIndex int
	volume           float64
	logger           *log.Logger
}

func NewAudioPlayer(audioContext *audio.Context, dirName string, logger *log.Logger) (*AudioPlayer, error) {
	files, err := os.ReadDir(dirName)
	if err != nil {
		return nil, fmt.Errorf("failed to read music directory: %v", err)
	}

	var allMusicFiles []*os.File
	for _, file := range files {
		musicFile, err := os.Open("music/" + file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to open music: %v", err)
		}
		allMusicFiles = append(allMusicFiles, musicFile)
	}
	rand.Shuffle(len(allMusicFiles), func(i, j int) { allMusicFiles[i], allMusicFiles[j] = allMusicFiles[j], allMusicFiles[i] })

	audioPlayer := &AudioPlayer{
		audioContext:  audioContext,
		allMusicFiles: allMusicFiles,
		logger:        logger,
	}

	return audioPlayer, nil
}

func (audioPlayer *AudioPlayer) Update() error {
	if !audioPlayer.player.IsPlaying() {
		err := audioPlayer.Next()
		if err != nil {
			return err
		}
	}
	return nil
}

func (audioPlayer *AudioPlayer) Next() error {
	audioPlayer.currentSongIndex++
	if audioPlayer.currentSongIndex >= len(audioPlayer.allMusicFiles) {
		audioPlayer.currentSongIndex = 0
	}
	return audioPlayer.play()
}

func (audioPlayer *AudioPlayer) Before() error {
	audioPlayer.currentSongIndex--
	if audioPlayer.currentSongIndex < 0 {
		audioPlayer.currentSongIndex = len(audioPlayer.allMusicFiles) - 1
	}
	return audioPlayer.play()
}

func (audioPlayer *AudioPlayer) play() error {
	file := audioPlayer.allMusicFiles[audioPlayer.currentSongIndex]
	mp3Stream, err := mp3.DecodeF32(file)
	if err != nil {
		return fmt.Errorf("failed to decode music: %v", err)
	}

	if audioPlayer.player != nil {
		if err = audioPlayer.player.Close(); err != nil {
			return fmt.Errorf("failed to close audio player: %v", err)
		}
	}
	audioPlayer.logger.Printf("[INFO] Playing music file %s\n", file.Name())
	player, err := audioPlayer.audioContext.NewPlayerF32(mp3Stream)
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}
	audioPlayer.player = player
	audioPlayer.currentSongName = path.Base(file.Name())
	audioPlayer.player.Play()
	audioPlayer.player.SetVolume(audioPlayer.volume)
	return nil
}

func (audioPlayer *AudioPlayer) Play() {
	if audioPlayer.player == nil {
		_ = audioPlayer.Next()
	}
	audioPlayer.player.Play()
}

func (audioPlayer *AudioPlayer) Pause() {
	audioPlayer.player.Pause()
}

func (audioPlayer *AudioPlayer) SetVolume(volume float64) {
	audioPlayer.player.SetVolume(volume)
	audioPlayer.volume = volume
}

func (audioPlayer *AudioPlayer) Volume() float64 {
	return audioPlayer.player.Volume()
}

func (audioPlayer *AudioPlayer) SongName() string {
	return strings.TrimSuffix(audioPlayer.currentSongName, ".mp3")
}
