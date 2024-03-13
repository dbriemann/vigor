package vigor

import "github.com/hajimehoshi/ebiten/v2"

type visual struct {
	scale   Vec2[float32]
	visible bool
}

func newVisual() visual {
	v := visual{
		scale:   Vec2[float32]{X: 1, Y: 1},
		visible: true,
	}
	return v
}

func (v *visual) FlipX() {
	v.scale.X *= -1
}

func (v *visual) FlipY() {
	v.scale.Y *= -1
}

func (v *visual) Scale(x, y float32) {
	v.scale.X *= x
	v.scale.Y *= y
}

func (v *visual) Visible() bool {
	return v.visible
}

func (v *visual) Show(on bool) {
	v.visible = on
}

func (v *visual) transform(op *ebiten.DrawImageOptions, width, height int) {
	tx := 0.0
	ty := 0.0

	op.GeoM.Scale(float64(v.scale.X), float64(v.scale.Y))

	if v.scale.X < 0 {
		tx = float64(width) * -1.0 * float64(v.scale.X)
	}
	if v.scale.Y < 0 {
		ty = float64(height) * -1.0 * float64(v.scale.Y)
	}

	op.GeoM.Translate(tx, ty)

	// TODO: debug wireframe
	// vector.StrokeRect(target, s.pos.X, s.pos.Y, float32(s.dim.X), float32(s.dim.Y), 2, color.White, false)
}

// scale = 1 --> no translation
// scale = -1 --> translate by 1 x size
// scale =
// scale| trans
// -----|-----
//   1  |  0
//  -1  |  1
//   2  |  0
//  -2  |  2
