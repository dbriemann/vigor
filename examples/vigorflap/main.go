package main

// TODO: attribution

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"

	"github.com/dbriemann/vigor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 160
	screenHeight = 240
)

var (
	dt               = 1.0 / float32(ebiten.TPS())
	score     uint32 = 0
	highscore uint32 = 0
)

type XY struct {
	X float32
	Y float32
}

type Dove struct {
	pos           XY
	vel           XY
	accel         XY
	flip          bool
	width, height int
	activeAnim    *vigor.Animation
	animSail      *vigor.Animation
	animFlap      *vigor.Animation
}

func (d *Dove) Update() {

	// TODO: fix distance to border..
	if d.pos.X <= 0 {
		score++
		d.vel.X *= -1
		d.flip = false
	}
	if d.pos.X+float32(d.width) >= screenWidth {
		score++
		d.vel.X *= -1
		d.flip = true
	}

	d.pos.X += d.vel.X * dt
	d.pos.Y += d.vel.Y * dt

	d.vel.X += d.accel.X * dt
	d.vel.Y += d.accel.Y * dt

	if d.activeAnim == d.animFlap && d.activeAnim.Finished {
		d.activeAnim = d.animSail
		d.activeAnim.Reset()
		d.activeAnim.Run()
	}

	d.activeAnim.Update()
}

func (d *Dove) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	d.activeAnim.Draw(target, op)
}

type Paddle struct {
	y       float32
	targetY float32
}

type Game struct {
	man        vigor.ResourceManager
	background *ebiten.Image
	dove       Dove
	spikes     *ebiten.Image
}

func (g *Game) Init() {
	score = 0
	g.dove.pos.X = screenWidth / 2
	g.dove.pos.Y = screenHeight / 2
	g.dove.vel.X = 0
	g.dove.vel.Y = 0
	g.dove.accel.X = 0
	g.dove.accel.Y = 0
	g.dove.flip = false
	g.dove.activeAnim = g.dove.animSail
	g.dove.activeAnim.Run()
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.dove.accel.Y == 0 {
			g.dove.accel.Y = 500
			g.dove.vel.X = 80
		}
		g.dove.activeAnim = g.dove.animFlap
		g.dove.activeAnim.Reset()
		g.dove.activeAnim.Finished = false
		g.dove.activeAnim.Run()
		g.dove.vel.Y = -screenHeight
	}

	g.dove.Update()

	if g.dove.pos.Y <= float32(g.spikes.Bounds().Dy()) ||
		g.dove.pos.Y+float32(g.dove.width) >= float32(screenHeight-g.spikes.Bounds().Dy()) {
		g.Over()
		fmt.Println("high", highscore)
	}

	return nil
}

func (g *Game) Over() {
	if score > highscore {
		highscore = score
	}
	score = 0
	g.Init()
}

func (g *Game) Draw(target *ebiten.Image) {
	target.DrawImage(g.background, nil)

	topSpikesOp := &ebiten.DrawImageOptions{}
	topSpikesOp.GeoM.Scale(1, -1)
	topSpikesOp.GeoM.Translate(0, float64(g.spikes.Bounds().Dy()))
	target.DrawImage(g.spikes, topSpikesOp)

	bottomSpikesOp := &ebiten.DrawImageOptions{}
	bottomSpikesOp.GeoM.Translate(0, float64(screenHeight-g.spikes.Bounds().Dy()))
	target.DrawImage(g.spikes, bottomSpikesOp)

	doveOp := &ebiten.DrawImageOptions{}
	if g.dove.flip {
		doveOp.GeoM.Scale(-1, 1)
		doveOp.GeoM.Translate(float64(g.dove.width), 0)
	}
	doveOp.GeoM.Translate(float64(g.dove.pos.X), float64(g.dove.pos.Y))
	g.dove.Draw(target, doveOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{
		man: vigor.NewResourceManager(),
	}

	if err := g.man.LoadConfig("assets/config.json"); err != nil {
		panic(err)
	}

	// TODO: load or panic function in resource man
	bg, ok := g.man.Images["background"]
	if !ok {
		panic("no background image loaded")
	}
	g.background = bg

	spikes, ok := g.man.Images["spikes"]
	if !ok {
		panic("could not load spikes image")
	}
	g.spikes = spikes

	a, err := vigor.NewAnimation(g.man.AnimationTemplates["dove_sail"])
	if err != nil {
		panic(err)
	}
	g.dove.animSail = a

	a, err = vigor.NewAnimation(g.man.AnimationTemplates["dove_flap"])
	if err != nil {
		panic(err)
	}
	g.dove.animFlap = a

	g.Init()
	g.dove.width = g.dove.activeAnim.FrameWidth
	g.dove.height = g.dove.activeAnim.FrameHeight

	return g
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth*3, screenHeight*3)
	ebiten.SetWindowTitle("Vigorflap")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
