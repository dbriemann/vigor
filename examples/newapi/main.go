package main

import (
	"github.com/dbriemann/vigor"
)

const (
	screenWidth  = 160
	screenHeight = 240
)

type Game struct{}

func (g *Game) Init() {
}

func (g *Game) Update() {
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

type Knight struct {
	vigor.Sprite
}

func (k *Knight) Update() {
	k.Sprite.Update()
}

func NewKnight(x, y int) *Knight {
	k := &Knight{
		Sprite: *vigor.NewSprite("knight_attack1"),
	}
	k.SetPos(float32(x), float32(y))
	return k
}

func main() {
	g := Game{}

	vigor.InitGame()

	knight := NewKnight(50, 50)
	knight.Scale(3, 3)

	vigor.G.Add(knight)

	vigor.RunGame(&g)
}
