package vigor

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Init()
	Update()
	Layout(width, height int) (logicalWidth, logicalHeight int)
}

type internalGame struct {
	stage DisplayList
}

func (g *internalGame) Draw(target *ebiten.Image) {
	g.stage.draw(target)
}

func (g *internalGame) Update() error {
	g.stage.Update()
	G.externalGame.Update()
	return nil
}

func (g *internalGame) add(s Stageable) {
	g.stage.Add(s)
}

func (g *internalGame) Layout(width, height int) (logicalWidth, logicalHeight int) {
	return G.externalGame.Layout(width, height)
}

func InitGame() error {
	G.assets = NewAssetManager()

	if err := G.assets.LoadConfig(configFilePath); err != nil {
		return err
	}

	G.SetTPS(60)

	return nil
}

func RunGame(g Game) error {
	G.externalGame = g
	G.externalGame.Init()

	return ebiten.RunGame(&G.internalGame)
}
