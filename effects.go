package vigor

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

// TODO: somehow unify effect usage
type Effect interface {
	update(float32) bool
	draw(*ebiten.Image, *colorm.DrawImageOptions)

	Start()
	Stop()
	Reset()
}

type ShakeEffect struct {
	magnitudeX float32
	magnitudeY float32
	displaceX  float32
	displaceY  float32
	duration   float32
	runtime    float32
	running    bool
}

func NewShakeEffect(duration, magX, magY float32) *ShakeEffect {
	e := &ShakeEffect{
		duration:   duration,
		runtime:    0,
		magnitudeX: magX,
		magnitudeY: magY,
		running:    false,
	}

	return e
}

func (e *ShakeEffect) Update(dt float32) bool {
	if !e.running {
		return false
	}

	e.displaceX = rand.Float32()*e.magnitudeX - e.magnitudeX/2
	e.displaceY = rand.Float32()*e.magnitudeY - e.magnitudeY/2

	e.runtime += dt

	if e.runtime >= e.duration {
		e.runtime = e.duration
		e.running = false
	}

	return !e.running
}

// TODO: probably should be Draw() to have unique interface for all effects.
func (e *ShakeEffect) Apply(op *ebiten.DrawImageOptions) {
	if !e.running {
		return
	}
	op.GeoM.Translate(float64(e.displaceX), float64(e.displaceY))
}

func (e *ShakeEffect) Reset() {
	e.runtime = 0
	e.running = false
	e.Stop()
}

func (e *ShakeEffect) Start() {
	e.running = true
}

func (e *ShakeEffect) Stop() {
	e.running = false
}

type FlashEffect struct {
	overlay  *ebiten.Image
	tweenSeq *gween.Sequence
	value    float32
	finished bool
	running  bool
}

func NewFlashEffect(image *ebiten.Image, duration float32, in, out ease.TweenFunc) *FlashEffect {
	e := &FlashEffect{
		overlay:  ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy()),
		finished: false,
		running:  false,
		tweenSeq: gween.NewSequence(
			gween.New(0, 1, duration/2, in),
			gween.New(1, 0, duration/2, out),
		),
	}
	e.overlay.Fill(color.White)

	return e
}

func (e *FlashEffect) Update(dt float32) bool {
	if !e.running {
		return false
	}
	if e.finished {
		return true
	}

	val, _, finished := e.tweenSeq.Update(dt)
	e.value = val
	e.finished = finished
	return finished
}

func (e *FlashEffect) Draw(target *ebiten.Image, op *colorm.DrawImageOptions) {
	if !e.running || e.finished {
		return
	}
	cm := colorm.ColorM{}
	cm.ChangeHSV(1, 1, float64(e.value))
	colorm.DrawImage(target, e.overlay, cm, op)
}

func (e *FlashEffect) Reset() {
	e.tweenSeq.Reset()
	e.finished = false
	e.Stop()
}

func (e *FlashEffect) Start() {
	e.running = true
}

func (e *FlashEffect) Stop() {
	e.running = false
}
