package vigor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type stageable interface {
	draw(*ebiten.Image, colorm.DrawImageOptions)

	Id() uint64
	Update()
	Visible() bool
	Show(bool)
}

type effected interface {
	stageable
	Dim() *Vec2[uint32]
	ApplyEffect(e Effect)
}
