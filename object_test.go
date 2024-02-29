package vigor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestObject(x, y float32, w, h uint32) (e Object) {
	G.idcounter++
	e.id = G.createId()
	e.motionDisabled = false
	e.pos = Vec2[float32]{x, y}
	e.dim = Vec2[uint32]{w, h}
	return
}

func TestCollides(t *testing.T) {
	testcases := []struct {
		name     string
		obj1     Object
		obj2     Object
		collides bool
	}{
		{
			name:     "partial overlap: collision",
			obj1:     newTestObject(10, 10, 5, 5),
			obj2:     newTestObject(14, 10, 3, 20),
			collides: true,
		},
		{
			name:     "both objects are equal: collision",
			obj1:     newTestObject(10, 10, 5, 5),
			obj2:     newTestObject(10, 10, 5, 5),
			collides: true,
		},
		{
			name:     "one object has zero area: no collision",
			obj1:     newTestObject(1, 1, 0, 3),
			obj2:     newTestObject(1, 1, 5, 5),
			collides: false,
		},
		{
			name:     "objects do not overlap: no collision",
			obj1:     newTestObject(1, 1, 5, 5),
			obj2:     newTestObject(10, 10, 5, 5),
			collides: false,
		},
	}

	for _, tc := range testcases {
		result := Collides(&tc.obj1, &tc.obj2)
		assert.Equal(t, tc.collides, result)
	}
}

// TODO: Test collisions with movement.
