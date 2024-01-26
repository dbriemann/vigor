package vigor

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type FlashEffect struct {
	overlay  *ebiten.Image
	tweenSeq *gween.Sequence
	value    float32
	finished bool
	running  bool
}

func NewFlashEffect(target *ebiten.Image, duration time.Duration, in, out ease.TweenFunc) *FlashEffect {
	e := &FlashEffect{
		overlay:  ebiten.NewImage(target.Bounds().Dx(), target.Bounds().Dy()),
		finished: false,
		running:  false,
		tweenSeq: gween.NewSequence(
			gween.New(0, 1, float32(duration.Seconds())/2, in),
			gween.New(1, 0, float32(duration.Seconds())/2, in),
		),
	}
	e.overlay.Fill(color.White)

	return e
}

func (e *FlashEffect) Update(dt float32) bool {
	if !e.running {
		return false
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
