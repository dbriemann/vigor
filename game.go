package vigor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	ebinput "github.com/quasilyte/ebitengine-input"
)

type Game interface {
	Init()
	Update()
	Layout(width, height int) (logicalWidth, logicalHeight int)
}

type internalGame struct {
	stage DisplayList
	input ebinput.System
}

func (g *internalGame) Draw(target *ebiten.Image) {
	// target.Fill(color.RGBA{0xff, 0, 0, 0xff})
	g.stage.draw(target)
	ebitenutil.DebugPrint(target, G.debugMsg)
}

func (g *internalGame) Update() error {
	g.input.Update()
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

func InitGame(g Game) error {
	// TODO: put this all in RunGame?
	G.assets = NewAssetManager()

	if err := G.assets.LoadConfig(configFilePath); err != nil {
		return err
	}

	G.SetTPS(60)

	G.internalGame.input.Init(ebinput.SystemConfig{
		DevicesEnabled: ebinput.AnyDevice,
	})

	G.externalGame = g
	G.externalGame.Init()

	return nil
}

func RunGame() error {
	return ebiten.RunGame(&G.internalGame)
}
