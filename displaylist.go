package vigor

import "github.com/hajimehoshi/ebiten/v2"

type DisplayList struct {
	// TODO: use better data structure than slice for adding, removing and being ordered.
	staged  []Stageable
	visible bool
}

func (d *DisplayList) Add(s Stageable) {
	d.staged = append(d.staged, s)
}

func (d *DisplayList) Remove(s Stageable) {
	id := s.Id()
	for i := 0; i < len(d.staged); i++ {
		if d.staged[i].Id() == id {
			d.staged = append(d.staged[:i], d.staged[i+1:]...)
			return
		}
	}
}

func (d *DisplayList) draw(target *ebiten.Image) {
	for i := 0; i < len(d.staged); i++ {
		if d.staged[i].Visible() {
			d.staged[i].draw(target)
		}
	}
}

func (d *DisplayList) Update() {
	for i := 0; i < len(d.staged); i++ {
		d.staged[i].Update()
	}
}

func (d *DisplayList) Visible() bool {
	return d.visible
}

func (d *DisplayList) Show(v bool) {
	d.visible = v
}
