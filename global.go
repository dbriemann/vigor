package vigor

var (
	igame internalGame
)

func Dt() float32 {
	return igame.dt
}

func TPS() uint32 {
	return igame.tps
}
