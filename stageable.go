package vigor

import "github.com/hajimehoshi/ebiten/v2"

type Stageable interface {
	draw(*ebiten.Image)

	Id() uint64
	Update()
	Visible() bool
	Show(bool)
}
