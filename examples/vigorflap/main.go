package main

// TODO: attribution

import (
	_ "embed"
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

type XY struct {
	X float32
	Y float32
}

type Dove struct {
	Pos           XY
	Vel           XY
	Accel         XY
	Flip          bool
	Width, Height int
	ActiveAnim    *vigor.Animation
	AnimSail      *vigor.Animation
	AnimFlap      *vigor.Animation
}

func (d *Dove) Update() {
	dt := 1.0 / float32(ebiten.TPS())

	d.Pos.X += d.Vel.X * dt
	d.Pos.Y += d.Vel.Y * dt

	d.Vel.X += d.Accel.X * dt
	d.Vel.Y += d.Accel.Y * dt

	if d.Pos.X <= 5 {
		d.Vel.X *= -1
		d.Flip = false
	}
	if d.Pos.X+10 >= screenWidth {
		d.Vel.X *= -1
		d.Flip = true
	}

	if d.ActiveAnim == d.AnimFlap && d.ActiveAnim.Finished {
		d.ActiveAnim = d.AnimSail
		d.ActiveAnim.Reset()
		d.ActiveAnim.Run()
	}

	d.ActiveAnim.Update()
}

func (d *Dove) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	d.ActiveAnim.Draw(target, op)
}

type Game struct {
	man        vigor.ResourceManager
	background *ebiten.Image
	dove       Dove
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.dove.Accel.Y == 0 {
			g.dove.Accel.Y = 500
			g.dove.Vel.X = 80
		}
		g.dove.ActiveAnim = g.dove.AnimFlap
		g.dove.ActiveAnim.Reset()
		g.dove.ActiveAnim.Finished = false
		g.dove.ActiveAnim.Run()
		g.dove.Vel.Y = -240
	}

	g.dove.Update()

	return nil
}

func (g *Game) Draw(target *ebiten.Image) {
	target.DrawImage(g.background, nil)
	op := &ebiten.DrawImageOptions{}
	if g.dove.Flip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(g.dove.Width), 0)
	}
	op.GeoM.Translate(float64(g.dove.Pos.X), float64(g.dove.Pos.Y))
	g.dove.Draw(target, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{
		man: vigor.NewResourceManager(),
		dove: Dove{
			Pos: XY{
				X: screenWidth / 2,
				Y: screenHeight / 2,
			},
		},
	}

	if err := g.man.LoadConfig("assets/config.json"); err != nil {
		panic(err)
	}

	bg, ok := g.man.Images["background"]
	if !ok {
		panic("no background image loaded")
	}
	g.background = bg

	a, err := vigor.NewAnimation(g.man.AnimationTemplates["dove_sail"])
	if err != nil {
		panic(err)
	}
	g.dove.AnimSail = a

	a, err = vigor.NewAnimation(g.man.AnimationTemplates["dove_flap"])
	if err != nil {
		panic(err)
	}
	g.dove.AnimFlap = a

	g.dove.ActiveAnim = g.dove.AnimSail
	g.dove.ActiveAnim.Run()

	return g
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Vigorflap")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
