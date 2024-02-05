package vigor

type Entity struct {
	pos            Vec2[float32]
	lastPos        Vec2[float32]
	vel            Vec2[float32]
	accel          Vec2[float32]
	dim            Vec2[uint32]
	motionDisabled bool
}

func NewEntity(pos Vec2[float32], dim Vec2[uint32], vel, accel Vec2[float32]) (e Entity) {
	e.pos = pos
	e.dim = dim
	e.vel = vel
	e.accel = accel
	e.motionDisabled = false

	return
}

func (e *Entity) SetMotion(enabled bool) {
	e.motionDisabled = !enabled
}

func (e Entity) Pos() Vec2[int] {
	// TODO: is Round better than Floor for pixel perfect positions?
	return Vec2Floor[float32, int](e.pos)
}

func (e *Entity) Update() {

}
