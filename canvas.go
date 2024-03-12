package vigor

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Canvas is similar to a sprite, but can be altered and is not animated.
type Canvas struct {
	image   *ebiten.Image
	visible bool

	visual
	Object
}

// TODO: do we need the Image alias to Canvas at all?
type Image = Canvas

func NewImage(name string) *Image {
	i := &Image{
		Object: NewObject(),
		visual: NewVisual(),

		visible: true,
		image:   G.assets.GetImageOrPanic(name),
	}

	i.SetDim(uint32(i.image.Bounds().Dx()), uint32(i.image.Bounds().Dy()))
	return i
}

func NewCanvas(width, height int) *Canvas {
	c := &Canvas{
		Object: NewObject(),
		visual: NewVisual(),

		visible: true,
		image:   ebiten.NewImage(width, height),
	}
	c.SetDim(uint32(width), uint32(height))

	return c
}

func (c *Canvas) draw(target *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(c.scale.X), float64(c.scale.Y))
	op.GeoM.Translate(float64(c.PixelPos().X), float64(c.PixelPos().Y))
	target.DrawImage(c.image, op)
}

func (s *Canvas) Show(v bool) {
	s.visible = v
}

func (c *Canvas) Visible() bool {
	return c.visible
}

func (c *Canvas) DrawFilledRect(x, y, width, height float32, col color.Color, antialias bool) {
	vector.DrawFilledRect(c.image, x, y, width, height, col, antialias)
}

func (c *Canvas) DrawRect(x, y, width, height, strokeWidth float32, col color.Color, antialias bool) {
	vector.StrokeRect(c.image, x, y, width, height, strokeWidth, col, antialias)
}

// TODO: other shapes
