package vigor

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tanema/gween/ease"
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

	// TODO: how set dim for sprites? adjust with scaling?
	s.Object.dim.X = uint32(s.activeAnim.AnimationTemplate.FrameWidth)
	s.Object.dim.Y = uint32(s.activeAnim.AnimationTemplate.FrameHeight)

	s.activeAnim.Run()

	return s
}

func (s *Sprite) draw(target *ebiten.Image) {
	if !s.Visible() {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(s.scale.X), float64(s.scale.Y))
	op.GeoM.Translate(float64(s.Pos().X), float64(s.Pos().Y))
	s.activeAnim.Draw(target, op)
	// TODO: debug wireframe
	// vector.StrokeRect(target, s.pos.X, s.pos.Y, float32(s.dim.X), float32(s.dim.Y), 2, color.White, false)
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

// SetTweenFunc sets the easing function for the active animation.
func (s *Sprite) SetTweenFunc(f ease.TweenFunc) {
	s.activeAnim.SetTweenFunc(f)
	s.activeAnim.InitTween()
}

func (s *Sprite) SetDuration(dur time.Duration) {
	s.activeAnim.SetDuration(dur)
	s.activeAnim.InitTween()
}

// func (s *Sprite) ActiveAnimation() *Animation {
// 	return s.activeAnim
// }
