package vigor

type Object struct {
	pos            Vec2[float32]
	lastPos        Vec2[float32]
	vel            Vec2[float32]
	accel          Vec2[float32]
	dim            Vec2[uint32]
	motionDisabled bool
}

func NewObject(x, y float32, width, height uint32) (e Object) {
	e.SetPos(x, y)
	e.SetDim(width, height)
	e.motionDisabled = true

	return
}

func (e *Object) SetPos(x, y float32) {
	e.pos.X = x
	e.pos.Y = y
}

func (e *Object) SetVel(x, y float32) {
	e.vel.X = x
	e.vel.Y = y
	e.SetMotion(true)
}

func (e *Object) SetAccel(x, y float32) {
	e.accel.X = x
	e.accel.Y = y
	e.SetMotion(true)
}

func (e *Object) SetDim(x, y uint32) {
	e.dim.X = x
	e.dim.Y = y
}

func (e *Object) SetMotion(enabled bool) {
	e.motionDisabled = !enabled
}

func (e Object) Pos() Vec2[int] {
	// TODO: is Round better than Floor for pixel perfect positions?
	return Vec2Floor[float32, int](e.pos)
}

func (e *Object) Update() {

}

func Collides(obj1, obj2 *Object) bool {
	// TODO: use lastPos for better collision detection.

	// If one of the rectangles does have a zero area, there is no intersection.
	if obj1.dim.X <= 0 || obj2.dim.X <= 0 || obj1.dim.Y <= 0 || obj2.dim.Y <= 0 {
		return false
	}

	// If one of the rectangles is left of the other, there is no intersection.
	if obj1.pos.X > obj2.pos.X+float32(obj2.dim.X) || obj1.pos.X+float32(obj1.dim.X) < obj2.pos.X {
		return false
	}

	// If one of the rectangles is above the other, there is no intersection.
	if obj1.pos.Y > obj2.pos.Y+float32(obj2.dim.Y) || obj1.pos.Y+float32(obj1.dim.Y) < obj2.pos.Y {
		return false
	}

	return true
}
