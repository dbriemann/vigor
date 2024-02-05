package vigor

import "github.com/hajimehoshi/ebiten/v2"

type game interface {
	Init()
	Update()
	Layout(width, height int) (logicalWidth, logicalHeight int)
}

type internalGame struct {
	g game

	tps uint32
	dt  float32
}

func (g *internalGame) Draw(target *ebiten.Image) {
	// TODO: internal draw stuff
}

func (g *internalGame) SetTPS(tps uint32) {
	if tps >= 1 {
		g.tps = tps
		g.dt = 1.0 / float32(tps)
	}
}

func (g *internalGame) Update() error {
	// TODO: internal update stuff
	g.g.Update()
	return nil
}

func (g *internalGame) Layout(width, height int) (logicalWidth, logicalHeight int) {
	return g.g.Layout(width, height)
}

func RunGame(g game) error {
	igame.g = g

	// TODO: should we put any data into game? How separate it from G?

	// TODO: setup globals: states etc

	return ebiten.RunGame(&igame)
}
