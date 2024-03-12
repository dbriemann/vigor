package main

import (
	"image/color"

	"github.com/dbriemann/vigor"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	screenWidth      = 160
	screenHeight     = 240
	score        int = 0
	highscore    int = 0

	ActionFlap input.Action = iota
)

var keymap = input.Keymap{
	ActionFlap: {input.KeySpace, input.KeyGamepadX},
}

type Game struct {
	input        *input.Handler
	dove         *Dove
	bouncerLeft  *Bouncer
	bouncerRight *Bouncer
	spikesTop    *vigor.Image
	spikesBottom *vigor.Image
}

func (g *Game) Init() {
	g.input = vigor.NewInputHandler(0, keymap)

	g.spikesTop = vigor.NewImage("spikes")
	g.spikesTop.SetPos(0, 50)
	g.spikesTop.FlipY()
	vigor.G.Add(g.spikesTop)

	g.bouncerLeft = NewBouncer(0, 0)
	g.bouncerLeft.Init()

	g.bouncerRight = NewBouncer(screenWidth-4, 0)
	g.bouncerRight.Init()

	// Dove is added last so it is painted last.
	g.dove = NewDove()
	g.dove.Init()
}

func (g *Game) Update() {
	if g.input.ActionIsJustPressed(ActionFlap) {
		if g.dove.Object.Accel().Y == 0 {
			g.dove.SetAccel(0, 500)
			g.dove.Accel().Y = 500
			g.dove.Vel().X = 80
		}
		g.dove.SetAnimation("dove_flap")
		g.dove.Vel().Y = -screenHeight
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

type Bouncer struct {
	back  *vigor.Canvas
	front *vigor.Canvas
}

func NewBouncer(x, y float32) *Bouncer {
	b := &Bouncer{
		// TODO: using 2 canvases sucks for this
		back:  vigor.NewCanvas(4, screenHeight),
		front: vigor.NewCanvas(4, screenHeight),
	}
	b.back.SetPos(x, y)
	b.front.SetPos(x, y)

	return b
}

func (b *Bouncer) Init() {
	b.back.DrawFilledRect(0, 0, 4, screenHeight, color.RGBA{53, 53, 61, 240}, false)
	b.front.DrawFilledRect(1, 1, 2, screenHeight-2, color.RGBA{100, 106, 125, 225}, false)
	vigor.G.Add(b.back)
	vigor.G.Add(b.front)
}

type Dove struct {
	vigor.Sprite
}

func NewDove() *Dove {
	d := &Dove{
		Sprite: *vigor.NewSprite("dove_sail", "dove_flap"),
	}
	return d
}

func (d *Dove) Init() {
	d.SetPos(screenWidth/2, screenHeight/2)
	vigor.G.Add(d)
}

func (d *Dove) Update() {
	anim, _, finished := d.Animation()
	if anim == "dove_flap" && finished {
		d.ResetAnimation()
		d.SetAnimation("dove_sail")
	}
	d.Sprite.Update()
}

func main() {
	vigor.SetWindowSize(3*screenWidth, 3*screenHeight)
	g := Game{}

	vigor.InitGame(&g)

	vigor.RunGame()
}
