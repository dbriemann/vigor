package vigor

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type DisplayGroup struct {
	// TODO: use better data structure than slice for adding, removing and being ordered.
	staged  []Stageable
	visible bool
}

func (d *DisplayGroup) Add(s Stageable) {
	d.staged = append(d.staged, s)
}

func (d *DisplayGroup) Remove(s Stageable) {
	id := s.Id()
	for i := 0; i < len(d.staged); i++ {
		if d.staged[i].Id() == id {
			d.staged = append(d.staged[:i], d.staged[i+1:]...)
			return
		}
	}
}

func (d *DisplayGroup) draw(target *ebiten.Image) {
	for i := 0; i < len(d.staged); i++ {
		if d.staged[i].Visible() {
			d.staged[i].draw(target)
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
