package vigor

import "golang.org/x/exp/constraints"

// TODO: move utils to proper place.

type Number interface {
	constraints.Integer | constraints.Float
}

type XY[T Number] struct {
	X T
	Y T
}

type WH[T Number] struct {
	W T
	H T
}

type Rect[T Number] struct {
	Point XY[T]
	Dim   WH[T]
}

func (r Rect[T]) Intersects(other Rect[T]) bool {
	// If one of the rectangles does have a zero area, there is no intersection.
	if r.Dim.W <= 0 || other.Dim.W <= 0 || r.Dim.H <= 0 || other.Dim.H <= 0 {
		return false
	}

	// If one of the rectangles is left of the other, there is no intersection.
	if r.Point.X > other.Point.X+other.Dim.W || r.Point.X+r.Dim.W < other.Point.X {
		return false
	}

	// If one of the rectangles is above the other, there is no intersection.
	if r.Point.Y > other.Point.Y+other.Dim.H || r.Point.Y+r.Dim.H < other.Point.Y {
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
