package vigor

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Init()
	Update()
	Layout(width, height int) (logicalWidth, logicalHeight int)
}

type internalGame struct {
}

func (g *internalGame) Draw(target *ebiten.Image) {
	// TODO: internal draw stuff
}

func (g *internalGame) Update() error {
	// TODO: internal update stuff
	G.exGame.Update()
	return nil
}

func (g *internalGame) Layout(width, height int) (logicalWidth, logicalHeight int) {
	return G.exGame.Layout(width, height)
}

func RunGame(g Game) error {
	G.exGame = g

	G.exGame.Init()

	return ebiten.RunGame(&G.inGame)
}
