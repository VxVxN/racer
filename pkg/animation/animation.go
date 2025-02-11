package animation

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Animation struct {
	x, y           float64
	scaleX, scaleY float64
	frameOX        int
	frameOY        int
	frameWidth     int
	frameHeight    int
	frameCount     int
	currentFrame   float64
	isRepeatable   bool
	start          bool
	callbackDone   bool
	image          *ebiten.Image
	player         *audio.Player
	callback       func()
}

func NewAnimation(image *ebiten.Image, frameOX, frameOY, frameWidth, frameHeight, frameCount int) *Animation {
	return &Animation{
		frameOX:      frameOX,
		frameOY:      frameOY,
		frameWidth:   frameWidth,
		frameHeight:  frameHeight,
		frameCount:   frameCount,
		image:        image,
		isRepeatable: true,
		scaleX:       1.0,
		scaleY:       1.0,
	}
}

func (animation *Animation) Update(speed float64) {
	if !animation.start {
		return
	}

	if animation.player != nil {
		animation.player.Play()
	}
	if int(animation.currentFrame) >= animation.frameCount && !animation.callbackDone {
		animation.callback()
		animation.callbackDone = true
	}
	animation.currentFrame += speed
}

func (animation *Animation) Draw(screen *ebiten.Image) {
	if !animation.start {
		return
	}
	if !animation.isRepeatable && int(animation.currentFrame) >= animation.frameCount {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(animation.x, animation.y)
	op.GeoM.Scale(animation.scaleX, animation.scaleY)
	i := int(animation.currentFrame)
	sx, sy := animation.frameOX+i*animation.frameWidth, animation.frameOY
	screen.DrawImage(animation.image.SubImage(image.Rect(sx, sy, sx+animation.frameWidth, sy+animation.frameHeight)).(*ebiten.Image), op)
}

func (animation *Animation) SetPosition(x, y float64) {
	animation.x = x
	animation.y = y
}

func (animation *Animation) SetRepeatable(enabled bool) {
	animation.isRepeatable = enabled
}

func (animation *Animation) Start() {
	animation.start = true
}

func (animation *Animation) Stop() {
	animation.start = false
}

func (animation *Animation) Reset() {
	animation.currentFrame = 0
	animation.callbackDone = false
	animation.start = false
	if animation.player != nil {
		animation.player.Pause()
		animation.player.SetPosition(0)
	}
}

func (animation *Animation) SetCallback(callback func()) {
	animation.callback = callback
}

func (animation *Animation) SetScale(scaleX, scaleY float64) {
	animation.scaleX = scaleX
	animation.scaleY = scaleY
}

func (animation *Animation) SetSound(audioContext *audio.Context, fileName string) error {
	musicFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open music: %v", err)
	}
	mp3Stream, err := mp3.DecodeF32(musicFile)
	if err != nil {
		return fmt.Errorf("failed to decode music: %v", err)
	}

	player, err := audioContext.NewPlayerF32(mp3Stream)
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}
	animation.player = player

	return nil
}

func (animation *Animation) SetVolume(volume float64) {
	animation.player.SetVolume(volume)
}
