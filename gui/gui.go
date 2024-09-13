package gui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button struct {
	x, y          int
	width, height int
}

func NewButton(x, y, width, height int) *Button {
	return &Button{x: x, y: y, width: width, height: height}
}

func (b *Button) IsClicked() bool {
	result := false
	mouseX, mouseY := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if mouseX >= b.x && mouseX <= b.x+b.width && mouseY >= b.y && mouseY <= b.y+b.height {
			result = true
		}
	}
	return result
}

func (b *Button) Draw(screen *ebiten.Image) {
	image := ebiten.NewImage(b.width, b.height)
	color := color.RGBA{255, 255, 255, 255}
	image.Fill(color)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(image, op)
}
