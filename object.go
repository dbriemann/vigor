package vigor

// Object represents any entity that has a position, updates. It can either be with or without motion.
type Object struct {
	pos            Vec2[float32]
	lastPos        Vec2[float32]
	vel            Vec2[float32]
	accel          Vec2[float32]
	dim            Vec2[uint32]
	id             uint64
	motionDisabled bool
}

func NewObject() (e Object) {
	G.idcounter++
	e.id = G.createId()
	e.motionDisabled = true
	return
}

func (e *Object) Id() uint64 {
	return e.id
}

func (e *Object) SetPos(x, y float32) {
	e.lastPos.X = x
	e.lastPos.Y = y
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

func (e *Object) PixelPos() Vec2[int] {
	// TODO: is Round better than Floor for pixel perfect positions?
	return Vec2Floor[float32, int](e.pos)
}

func (e *Object) Pos() *Vec2[float32] {
	return &e.pos
}

func (e *Object) Vel() *Vec2[float32] {
	return &e.vel
}

func (e *Object) Accel() *Vec2[float32] {
	return &e.accel
}

func (e *Object) Dim() *Vec2[uint32] {
	return &e.dim
}

func (e *Object) Update() {
	if e.motionDisabled {
		return
	}

	// TODO: angular velocity

	e.pos.X += e.vel.X * G.Dt()
	e.pos.Y += e.vel.Y * G.Dt()

	e.vel.X += e.accel.X * G.Dt()
	e.vel.Y += e.accel.Y * G.Dt()
}

type positionable interface {
	Pos() *Vec2[float32]
	Dim() *Vec2[uint32]
}

func Collides(obj1, obj2 positionable) bool {
	// TODO: use lastPos for better collision detection.

	// If one of the rectangles does have a zero area, there is no intersection.
	if obj1.Dim().X <= 0 || obj2.Dim().X <= 0 || obj1.Dim().Y <= 0 || obj2.Dim().Y <= 0 {
		return false
	}

	// If one of the rectangles is left of the other, there is no intersection.
	if obj1.Pos().X > obj2.Pos().X+float32(obj2.Dim().X) || obj1.Pos().X+float32(obj1.Dim().X) < obj2.Pos().X {
		return false
	}

	// If one of the rectangles is above the other, there is no intersection.
	if obj1.Pos().Y > obj2.Pos().Y+float32(obj2.Dim().Y) || obj1.Pos().Y+float32(obj1.Dim().Y) < obj2.Pos().Y {
		return false
	}

	return true
}
