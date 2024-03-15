package vigor

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Vec2[T Number] struct {
	X T
	Y T
}

func (v Vec2[T]) Normalize() (out Vec2[T]) {
	if v.X == 0 && v.Y == 0 {
		return v
	}
	len := math.Sqrt(float64(v.X*v.X + v.Y + v.Y))
	out.X = T(float64(v.X) / len)
	out.Y = T(float64(v.Y) / len)
	return
}

func (v Vec2[T]) Rotate(rad float64) (out Vec2[T]) {
	out.X = T(math.Cos(rad))*v.X - T(math.Sin(rad))*v.Y
	out.Y = T(math.Sin(rad))*v.X - T(math.Cos(rad))*v.Y
	return
}

func (v Vec2[T]) Multiply(s T) (out Vec2[T]) {
	out.X = v.X * s
	out.Y = v.Y * s
	return
}

func Vec2ToType[T, U Number](in Vec2[T]) (out Vec2[U]) {
	out.X = U(in.X)
	out.Y = U(in.Y)
	return
}

func Vec2Floor[T Number, U constraints.Integer](in Vec2[T]) (out Vec2[U]) {
	out.X = U(math.Floor(float64(in.X)))
	out.Y = U(math.Floor(float64(in.Y)))
	return
}

type Rect[T Number] struct {
	Point Vec2[T]
	Dim   Vec2[T]
}

func (r Rect[T]) Intersects(other Rect[T]) bool {
	// If one of the rectangles does have a zero area, there is no intersection.
	if r.Dim.X <= 0 || other.Dim.X <= 0 || r.Dim.Y <= 0 || other.Dim.Y <= 0 {
		return false
	}

	// If one of the rectangles is left of the other, there is no intersection.
	if r.Point.X > other.Point.X+other.Dim.X || r.Point.X+r.Dim.X < other.Point.X {
		return false
	}

	// If one of the rectangles is above the other, there is no intersection.
	if r.Point.Y > other.Point.Y+other.Dim.Y || r.Point.Y+r.Dim.Y < other.Point.Y {
		return false
	}

	return true
}

// func Intersects[T Number](r, r2 Rect[T]) bool {
// 	// If one of the rectangles does have a zero area, there is no intersection.
// 	if r.Dim.W <= 0 || r2.Dim.W <= 0 || r.Dim.H <= 0 || r2.Dim.H <= 0 {
// 		return false
// 	}

// 	// If one of the rectangles is left of the other, there is no intersection.
// 	if r.Point.X > r2.Point.X+r2.Dim.W || r.Point.X+r.Dim.W < r2.Point.X {
// 		return false
// 	}

// 	// If one of the rectangles is above the other, there is no intersection.
// 	if r.Point.Y > r2.Point.Y+r2.Dim.H || r.Point.Y+r.Dim.H < r2.Point.Y {
// 		return false
// 	}

// 	return true
// }
