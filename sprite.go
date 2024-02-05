package vigor

import "github.com/hajimehoshi/ebiten/v2"

// Sprite represents every entity that has a position, is updated and drawn.
type Sprite struct {
	img       *ebiten.Image
	animated  bool
	animation *Animation
}
