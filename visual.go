package vigor

type visual struct {
	scale Vec2[float32]
	trans Vec2[float32]
}

func NewVisual() visual {
	v := visual{
		scale: Vec2[float32]{X: 1, Y: 1},
		trans: Vec2[float32]{X: 0, Y: 0},
	}
	return v
}

func (v *visual) FlipX() {
	v.scale.X *= -1
	// TODO: translate on flips?
	// CONTINUE HERE
}

func (v *visual) FlipY() {
	v.scale.Y *= -1
}

func (v *visual) Scale(x, y float32) {
	v.scale.X *= x
	v.scale.Y *= y
}
