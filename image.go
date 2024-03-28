package vigor

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var _ effected = (*Image)(nil)

// Image is similar to a sprite, but can be altered and is not animated.
type Image struct {
	effects []Effect
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
		effects: img.effects,
	}
	return i
}

func NewImage(name string) *Image {
	i := &Image{
		Object:  NewObject(),
		visual:  newVisual(),
		effects: []Effect{},

		visible: true,
		image:   G.assets.GetImageOrPanic(name),
	}

	i.SetDim(uint32(i.image.Bounds().Dx()), uint32(i.image.Bounds().Dy()))
	return i
}

func NewCanvas(width, height int) *Image {
	c := &Image{
		Object:  NewObject(),
		visual:  newVisual(),
		effects: []Effect{},

		visible: true,
		image:   ebiten.NewImage(width, height),
	}
	c.SetDim(uint32(width), uint32(height))

	return c
}

func (i *Image) ApplyEffect(e Effect) {
	e.Reset()
	e.Start()
	i.effects = append(i.effects, e)
}

func (i *Image) Update() {
	i.Object.Update()
	for j := 0; j < len(i.effects); j++ {
		finished := i.effects[j].Update()
		if finished {
			i.effects = append(i.effects[:j], i.effects[j+1:]...)
		}
	}
}

func (c *Image) draw(target *ebiten.Image, op colorm.DrawImageOptions) {
	cm := colorm.ColorM{}
	c.transform(&op, int(c.Dim().X), int(c.Dim().Y))
	op.GeoM.Translate(float64(c.PixelPos().X), float64(c.PixelPos().Y))
	for i := 0; i < len(c.effects); i++ {
		c.effects[i].modifyDraw(&op)
	}
	colorm.DrawImage(target, c.image, cm, &op)
	for i := 0; i < len(c.effects); i++ {
		c.effects[i].draw(target, op)
	}
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
