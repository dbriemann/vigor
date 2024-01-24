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

func Intersects[T Number](r1, r2 Rect[T]) bool {
	// If one of the rectangles does have a zero area, there is no intersection.
	if r1.Dim.W <= 0 || r2.Dim.W <= 0 || r1.Dim.H <= 0 || r2.Dim.H <= 0 {
		return false
	}

	// If one of the rectangles is left of the other, there is no intersection.
	if r1.Point.X > r2.Point.X+r2.Dim.W || r1.Point.X+r1.Dim.W < r2.Point.X {
		return false
	}

	// If one of the rectangles is above the other, there is no intersection.
	if r1.Point.Y > r2.Point.Y+r2.Dim.H || r1.Point.Y+r1.Dim.H < r2.Point.Y {
		return false
	}

	return true
}
