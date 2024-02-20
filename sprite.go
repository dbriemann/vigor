package vigor

import "github.com/hajimehoshi/ebiten/v2"

// Sprite represents every entity that has a position, is updated and drawn.
type Sprite struct {
	Object

	img       *ebiten.Image
	animated  bool
	animation *Animation
}

func NewSprite(imageName string) *Sprite {
	s := &Sprite{}

	return s
}

func (s *Sprite) draw(target *ebiten.Image) {

}

func (e *Sprite) Update() {

}

// TODO: continue here
