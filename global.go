package vigor

// TODO: make glob thread safe.
var G glob

type glob struct {
	internalGame internalGame
	externalGame Game
	assets       AssetManager
	tps          uint32
	dt           float32
	idcounter    uint64
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

func (g *glob) Add(s Stageable) {
	g.internalGame.add(s)
}

func SetConfigFile(cfgFilePath string) {
	configFilePath = cfgFilePath
}
