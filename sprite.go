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
	activeAnim     *Animation
	animations     map[string]*Animation
	activeAnimName string
	scale          Vec2[float32]
	visible        bool

	Object
}

// NewSprite takes any amount of animations by name. These animations must exist in the asset manager.
// The first animation is used as default animation.
func NewSprite(animNames ...string) *Sprite {
	s := &Sprite{
		Object:     NewObject(),
		animations: map[string]*Animation{},
		scale:      Vec2[float32]{X: 1, Y: 1},
		visible:    true,
	}
	for i, name := range animNames {
		anim, err := NewAnimation(G.assets.GetAnimTemplateOrPanic(name))
		if err == nil {
			s.animations[name] = anim
			if i == 0 {
				s.activeAnim = anim
				s.activeAnimName = name
			}
		}
		// For any error we just skip creating the animation for now.
		// TODO: logging / behavior ?
	}

	// TODO: how set dim/bbox for sprites? adjust with scaling?
	s.Object.dim.X = uint32(s.activeAnim.AnimationTemplate.FrameWidth)
	s.Object.dim.Y = uint32(s.activeAnim.AnimationTemplate.FrameHeight)

	s.activeAnim.Run()

	return s
}

func (s *Sprite) SetAnimation(name string) {
	anim, ok := s.animations[name]
	if !ok {
		// TODO: ??
		return
	}
	s.activeAnim.Stop()
	s.activeAnim = anim
	s.activeAnimName = name
	s.activeAnim.Run()
}

func (s *Sprite) Animation() (name string, paused, finished bool) {
	name = s.activeAnimName
	paused = s.activeAnim.Paused
	finished = s.activeAnim.Finished
	return
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

func (s *Sprite) StopAnimation() {
	s.activeAnim.Stop()
}

func (s *Sprite) StartAnimation() {
	s.activeAnim.Run()
}

func (s *Sprite) ResetAnimation() {
	s.activeAnim.Reset()
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
