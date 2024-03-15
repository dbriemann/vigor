package vigor

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Image is similar to a sprite, but can be altered and is not animated.
type Image struct {
	image   *ebiten.Image
	visible bool

	visual
	Object
}

func CopyImage(img *Image) *Image {
	i := &Image{
		image:   img.image,
		visible: img.visible,
		visual:  img.visual,
		Object:  img.Object,
	}
	return i
}

func NewImage(name string) *Image {
	i := &Image{
		Object: NewObject(),
		visual: newVisual(),

		visible: true,
		image:   G.assets.GetImageOrPanic(name),
	}

	i.SetDim(uint32(i.image.Bounds().Dx()), uint32(i.image.Bounds().Dy()))
	return i
}

func NewCanvas(width, height int) *Image {
	c := &Image{
		Object: NewObject(),
		visual: newVisual(),

		visible: true,
		image:   ebiten.NewImage(width, height), // TODO: should this be a named image in asset manager?
	}
	c.SetDim(uint32(width), uint32(height))

	return c
}

func (c *Image) draw(target *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	c.transform(op, int(c.Dim().X), int(c.Dim().Y))
	op.GeoM.Translate(float64(c.PixelPos().X), float64(c.PixelPos().Y))
	target.DrawImage(c.image, op)
}

func (s *Image) Show(v bool) {
	s.visible = v
}

func (c *Image) Visible() bool {
	return c.visible
}

func (c *Image) DrawFilledRect(x, y, width, height float32, col color.Color, antialias bool) {
	vector.DrawFilledRect(c.image, x, y, width, height, col, antialias)
}

func (c *Image) DrawRect(x, y, width, height, strokeWidth float32, col color.Color, antialias bool) {
	vector.StrokeRect(c.image, x, y, width, height, strokeWidth, col, antialias)
}

// TODO: other shapes
