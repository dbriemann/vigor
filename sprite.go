package vigor

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Stageable interface {
	draw(*ebiten.Image)

	Id() uint64
	Update()
	Visible() bool
	Show(bool)
}

// Sprite represents every entity that has a position, is updated and drawn.
type Sprite struct {
	activeAnim *Animation
	animations map[string]*Animation
	scale      Vec2[float32]
	visible    bool

	Object
}

func NewSprite(animNames ...string) *Sprite {
	s := &Sprite{
		Object:     NewObject(),
		animations: map[string]*Animation{},
		scale:      Vec2[float32]{X: 1, Y: 1},
		visible:    true,
	}
	for _, name := range animNames {
		anim, err := NewAnimation(G.assets.GetAnimTemplateOrPanic(name))
		if err == nil {
			s.animations[name] = anim
			s.activeAnim = anim // TODO:
		}
		// For any error we just skip creating the animation for now.
		// TODO: logging / behavior ?
	}

	s.activeAnim.Run()

	return s
}

func (s *Sprite) draw(target *ebiten.Image) {
	if !s.Visible() {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(s.scale.X), float64(s.scale.Y))
	s.activeAnim.Draw(target, op)
}

func (s *Sprite) Scale(x, y float32) {
	s.scale.X = x
	s.scale.Y = y
}

func (s *Sprite) Update() {
	s.activeAnim.Update(G.Dt())
	s.Object.Update()
}

func (s *Sprite) Visible() bool {
	return s.visible
}

func (s *Sprite) Show(v bool) {
	s.visible = v
}
