package vigor

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	ebinput "github.com/quasilyte/ebitengine-input"
)

// TODO: make glob thread safe.
var G glob

type glob struct {
	internalGame internalGame
	externalGame Game
	assets       AssetManager
	tps          uint32
	dt           float32
	idcounter    uint64
	debugMsg     string // HACK:
}

func (g *glob) createId() uint64 {
	g.idcounter++
	return g.idcounter
}

func (g *glob) Dt() float32 {
	return g.dt
}

func (g *glob) TPS() uint32 {
	return g.tps
}

func (g *glob) SetTPS(tps uint32) {
	if tps >= 1 {
		g.tps = tps
		g.dt = 1.0 / float32(tps)
	}
}

func (g *glob) Add(s stageable) {
	g.internalGame.add(s)
}

func (g *glob) ApplyEffect(e Effect) {
	e.Start()
	g.internalGame.effects = append(g.internalGame.effects, e)
}

func (g *glob) Remove(s stageable) {
	// TODO:
}

func SetConfigFile(cfgFilePath string) {
	configFilePath = cfgFilePath
}

func SetWindowSize(w, h int) {
	ebiten.SetWindowSize(w, h)
}

func NewInputHandler(id uint8, keymap ebinput.Keymap) *ebinput.Handler {
	return G.internalGame.input.NewHandler(id, keymap)
}

// HACK: this is to change.
func DebugPrintf(format string, a ...any) {
	G.debugMsg = fmt.Sprintf(format, a...)
}
