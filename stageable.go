package vigor

import "github.com/hajimehoshi/ebiten/v2"

type stageable interface {
	draw(*ebiten.Image)

	Id() uint64
	Update()
	Visible() bool
	Show(bool)
}

type effectable interface {
	stageable
	Dim() *Vec2[uint32]
	ApplyEffect(e Effect)
}
