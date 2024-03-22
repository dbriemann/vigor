package vigor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// TODO: ?
// var _ stageable = (*DisplayGroup)(nil)

type DisplayGroup struct {
	// TODO: use better data structure than slice for adding, removing and being ordered.
	staged  []stageable
	visible bool
}

func (d *DisplayGroup) Add(s stageable) {
	d.staged = append(d.staged, s)
}

func (d *DisplayGroup) Remove(s stageable) {
	id := s.Id()
	for i := 0; i < len(d.staged); i++ {
		if d.staged[i].Id() == id {
			d.staged = append(d.staged[:i], d.staged[i+1:]...)
			return
		}
	}
}

func (d *DisplayGroup) draw(target *ebiten.Image, op colorm.DrawImageOptions) {
	for i := 0; i < len(d.staged); i++ {
		if d.staged[i].Visible() {
			d.staged[i].draw(target, op)
		}
	}
}

func (d *DisplayGroup) Update() {
	for i := 0; i < len(d.staged); i++ {
		d.staged[i].Update()
	}
}

func (d *DisplayGroup) Visible() bool {
	return d.visible
}

func (d *DisplayGroup) Show(v bool) {
	d.visible = v
}
