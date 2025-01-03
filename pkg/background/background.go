package background

import "github.com/hajimehoshi/ebiten/v2"

type Background struct {
	screenWidth float64
	image       *ebiten.Image
	y           float64
}

func New(image *ebiten.Image, screenWidth float64) *Background {
	return &Background{
		image:       image,
		screenWidth: screenWidth,
	}
}

func (background *Background) Update(scrollSpeed float64) {
	background.y += scrollSpeed
	if background.y > float64(background.image.Bounds().Dy()) {
		background.y = 0
	}
}

func (background *Background) Draw(screen *ebiten.Image) {
	for i := -1; i < 5; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(background.screenWidth/2-float64(background.image.Bounds().Dx())/2, float64(background.image.Bounds().Dy()*i)+background.y)
		screen.DrawImage(background.image, op)
	}
}
