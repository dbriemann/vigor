package vigor

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

// Object represents any entity that has a position, updates. It can either be with or without motion.
type Object struct {
	tweenX         *gween.Tween
	tweenY         *gween.Tween
	pos            Vec2[float32]
	lastPos        Vec2[float32]
	vel            Vec2[float32]
	accel          Vec2[float32]
	dim            Vec2[uint32]
	id             uint64
	motionDisabled bool
}

func NewObject() (o Object) {
	G.idcounter++
	o.id = G.createId()
	o.motionDisabled = true
	return
}

func (o *Object) Id() uint64 {
	return o.id
}

func (o *Object) TweenTo(x, y, duration float32, f ease.TweenFunc) {
	o.SetMotion(true)
	o.tweenX = gween.New(o.pos.X, x, duration, f)
	o.tweenY = gween.New(o.pos.Y, y, duration, f)
}

func (o *Object) SetPos(x, y float32) {
	o.lastPos.X = x
	o.lastPos.Y = y
	o.pos.X = x
	o.pos.Y = y
}

func (o *Object) SetVel(x, y float32) {
	o.vel.X = x
	o.vel.Y = y
	o.SetMotion(true)
}

func (o *Object) SetAccel(x, y float32) {
	o.accel.X = x
	o.accel.Y = y
	o.SetMotion(true)
}

func (o *Object) SetDim(x, y uint32) {
	o.dim.X = x
	o.dim.Y = y
}

func (o *Object) SetMotion(enabled bool) {
	o.motionDisabled = !enabled
}

func (o *Object) PixelPos() Vec2[int] {
	// TODO: is Round better than Floor for pixel perfect positions?
	return Vec2Floor[float32, int](o.pos)
}

func (o *Object) Pos() *Vec2[float32] {
	return &o.pos
}

func (o *Object) Vel() *Vec2[float32] {
	return &o.vel
}

func (o *Object) Accel() *Vec2[float32] {
	return &o.accel
}

func (o *Object) Dim() *Vec2[uint32] {
	return &o.dim
}

func (o *Object) Update() {
	if o.motionDisabled {
		return
	}

	isTweening := false
	dt := G.Dt()
	if o.tweenX != nil {
		newx, finishedx := o.tweenX.Update(dt)
		o.pos.X = newx
		if finishedx {
			o.tweenX = nil
		}
	}
	if o.tweenY != nil {
		newy, finishedy := o.tweenY.Update(dt)
		o.pos.Y = newy
		if finishedy {
			o.tweenY = nil
		}
	}

	// If tweening is active we do not do the normal update routine.
	if isTweening {
		return
	}

	// TODO: angular velocity

	o.pos.X += o.vel.X * G.Dt()
	o.pos.Y += o.vel.Y * G.Dt()

	o.vel.X += o.accel.X * G.Dt()
	o.vel.Y += o.accel.Y * G.Dt()
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
