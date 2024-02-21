package main

import "github.com/dbriemann/vigor"

type Game struct {
}

func (g *Game) Init() {
}

func (g *Game) Update() {
}

func (g *Game) Layout(w, h int) (int, int) {

	return 0, 0
}

func main() {
	g := Game{}

	vigor.InitGame()

	vigor.RunGame(&g)
}
