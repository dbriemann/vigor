package vigor

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Stageable interface {
	draw(*ebiten.Image)

	Id() uint64
	Update()
	Visible() bool
}

// Sprite represents every entity that has a position, is updated and drawn.
type Sprite struct {
	animations map[string]*Animation

	Object
}

func NewSprite(animNames ...string) *Sprite {
	s := &Sprite{
		Object:     NewObject(),
		animations: map[string]*Animation{},
	}
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

func (e *Sprite) Visible() bool {
	return true
}
