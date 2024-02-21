package vigor

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite represents every entity that has a position, is updated and drawn.
type Sprite struct {
	Object

	img        *ebiten.Image
	animated   bool
	animations map[string]*Animation
}

func NewSprite(imageName string, animNames ...string) *Sprite {
	s := &Sprite{
		img: G.assets.GetImageOrPanic(imageName),
	}
	s.animated = len(animNames) == 0
	for _, name := range animNames {
		anim, err := NewAnimation(G.assets.GetAnimTemplateOrPanic(name))
		if err == nil {
			s.animations[name] = anim
		}
		// For any error we just skip creating the animation for now.
		// TODO: logging / behavior ?
	}

	return s
}

func (s *Sprite) draw(target *ebiten.Image) {

}

func (e *Sprite) Update() {
	e.Object.Update()
}

// TODO: continue here
