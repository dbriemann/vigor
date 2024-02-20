package vigor

var (
	// TODO: make glob thread safe.
	G glob
)

type glob struct {
	inGame internalGame
	exGame Game

	tps uint32
	dt  float32
}

func (g *glob) Dt() float32 {
	return g.dt
}

func (g *glob) TPS() uint32 {
	return g.tps
}

func (g *glob) SetTPS(tps uint32) {
	if g.tps >= 1 {
		g.tps = tps
		g.dt = 1.0 / float32(tps)
	}
}
