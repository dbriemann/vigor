package main

import (
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
	input *input.Handler
	dove  *Dove
}

func (g *Game) Init() {
	g.input = vigor.NewInputHandler(0, keymap)
	g.dove = NewDove()
	vigor.G.Add(g.dove)
}

func (g *Game) Update() {
	if g.input.ActionIsJustPressed(ActionFlap) {
		// TODO: make dove flap
		g.dove.SetAnimation("dove_flap")
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
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

func (d *Dove) Update() {
	anim, _, finished := d.Animation()
	if anim == "dove_flap" && finished {
		d.ResetAnimation()
		d.SetAnimation("dove_sail")
	}
	d.Sprite.Update()
}

func main() {
	vigor.SetWindowSize(4*screenWidth, 4*screenHeight)
	g := Game{}

	vigor.InitGame(&g)

	vigor.RunGame()
}
