package vigor

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Init()
	Update()
	Layout(width, height int) (logicalWidth, logicalHeight int)
}

type internalGame struct {
	// TODO: add stage NEXT NEXT NEXT
}

func (g *internalGame) Draw(target *ebiten.Image) {
	// TODO: internal draw stuff
}

func (g *internalGame) Update() error {
	// TODO: internal update stuff
	G.externalGame.Update()
	return nil
}

func (g *internalGame) Layout(width, height int) (logicalWidth, logicalHeight int) {
	return G.externalGame.Layout(width, height)
}

func InitGame() error {
	G.assets = NewAssetManager()

	if err := G.assets.LoadConfig(configFilePath); err != nil {
		return err
	}

	return nil
}

func RunGame(g Game) error {
	G.externalGame = g
	G.externalGame.Init()

	return ebiten.RunGame(&G.internalGame)
}
