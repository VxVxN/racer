package rectangle

type Rectangle struct {
	X, Y, Width, Height float64
}

func New(x, y, width, height float64) *Rectangle {
	return &Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (rectangle *Rectangle) Collision(otherRect *Rectangle) bool {
	return rectangle.X < otherRect.X+otherRect.Width &&
		rectangle.X+rectangle.Width > otherRect.X &&
		rectangle.Y < otherRect.Y+otherRect.Height &&
		rectangle.Y+rectangle.Height > otherRect.Y
}
