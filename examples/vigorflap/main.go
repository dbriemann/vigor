package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/dbriemann/vigor"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/tanema/gween/ease"
)

const (
	screenWidth  = 160
	screenHeight = 240

	ActionFlap input.Action = iota
)

var (
	score      int = 0
	highscore  int = 0
	paddleMinY int = 0
	paddleMaxY int = 0
)

var keymap = input.Keymap{
	ActionFlap: {input.KeySpace, input.KeyGamepadX},
}

type Game struct {
	input          *input.Handler
	dove           *Dove
	bouncerLeft    *Bouncer
	bouncerRight   *Bouncer
	paddleLeft     *Paddle
	paddleRight    *Paddle
	spikesTop      *vigor.Image
	spikesBottom   *vigor.Image
	featherEmitter *vigor.Emitter
	background     *vigor.Image
	flash          *vigor.FlashEffect
	shake          *vigor.ShakeEffect
	gameOverScene  bool
}

func (g *Game) Init() {
	score = 0

	g.input = vigor.NewInputHandler(0, keymap)

	g.background = vigor.NewImage("background")
	vigor.G.Add(g.background)
	g.background.SetPos(0, 0)
	g.flash = vigor.NewFlashEffect(g.background, 0.3, ease.Linear, ease.Linear)
	g.shake = vigor.NewShakeEffect(0.6, 6, 6)

	g.spikesTop = vigor.NewImage("spikes")
	g.spikesTop.SetPos(0, 0)
	g.spikesTop.FlipY()
	vigor.G.Add(g.spikesTop)

	spikesHeight := g.spikesTop.Dim().Y

	g.spikesBottom = vigor.NewImage("spikes")
	g.spikesBottom.SetPos(0, float32(screenHeight-spikesHeight))
	vigor.G.Add(g.spikesBottom)

	bouncerHeight := int(screenHeight-2*spikesHeight) - 4
	bouncerTy := spikesHeight + 2
	g.bouncerLeft = NewBouncer(0, float32(bouncerTy), bouncerHeight)
	g.bouncerLeft.Init()

	g.bouncerRight = NewBouncer(screenWidth-4, float32(bouncerTy), bouncerHeight)
	g.bouncerRight.Init()

	g.paddleLeft = NewPaddle(5, false)
	vigor.G.Add(g.paddleLeft)
	paddleWidth := g.paddleLeft.Dim().X
	g.paddleRight = NewPaddle(screenWidth-5-int(paddleWidth), true)
	vigor.G.Add(g.paddleRight)

	paddleMinY = int(spikesHeight) + 2
	paddleMaxY = screenHeight - 2*int(spikesHeight) - 4 - int(g.paddleLeft.Dim().Y)

	// Dove is added last so it is painted last.
	g.dove = NewDove()
	g.dove.Init()

	feather := vigor.NewImage("feather")
	g.featherEmitter = vigor.NewParticleEmitter(*feather, screenWidth/2, screenHeight/2, 10, 0)
	vigor.G.Add(g.featherEmitter)
}

func (g *Game) Over() {
	if score > highscore {
		highscore = score
	}

	g.featherEmitter.SetOrigin(g.dove.Pos().X, g.dove.Pos().Y)
	g.dove.Die()
	vigor.G.ApplyEffect(g.flash)
	vigor.G.ApplyEffect(g.shake)
	g.featherEmitter.Show(true)
	g.featherEmitter.Burst()
	g.gameOverScene = true
}

func (g *Game) Update() {
	// TODO: create font (placeable) entity
	scoreStr := fmt.Sprintf("\n\n\n\n\n\n             %d", score)
	if highscore > 0 {
		scoreStr += fmt.Sprintf("\n       high: %d", highscore)
	}
	vigor.DebugPrintf(scoreStr)

	if g.gameOverScene {
		if g.featherEmitter.ActiveParticles() == 0 {
			g.gameOverScene = false
			g.dove.Live()
			g.paddleLeft.PlaceOffScreen()
			g.paddleRight.PlaceOffScreen()
			score = 0
		}
		return
	}

	if g.input.ActionIsJustPressed(ActionFlap) {
		if g.dove.Object.Accel().Y == 0 {
			g.dove.SetAccel(0, 500)
			g.dove.Accel().Y = 500
			g.dove.Vel().X = 80
		}
		g.dove.SetAnimation("dove_flap")
		g.dove.Vel().Y = -screenHeight
	}

	if vigor.Collides(g.dove, g.spikesTop) ||
		vigor.Collides(g.dove, g.spikesBottom) ||
		vigor.Collides(g.dove, g.paddleLeft) ||
		vigor.Collides(g.dove, g.paddleRight) {
		// death
		g.Over()
	} else if vigor.Collides(g.dove, g.bouncerLeft.back) {
		score++
		g.dove.Vel().X *= -1
		g.dove.FlipX()
		g.bouncerLeft.back.ApplyEffect(g.bouncerLeft.flash)
		g.paddleRight.PlaceRandomly()
	} else if vigor.Collides(g.dove, g.bouncerRight.back) {
		score++
		g.dove.Vel().X *= -1
		g.dove.FlipX()
		g.bouncerRight.back.ApplyEffect(g.bouncerRight.flash)
		g.paddleLeft.PlaceRandomly()
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

type Bouncer struct {
	back  *vigor.Image
	front *vigor.Image
	flash *vigor.FlashEffect
}

func NewBouncer(x, y float32, height int) *Bouncer {
	b := &Bouncer{
		back:  vigor.NewCanvas(4, height),
		front: vigor.NewCanvas(4, height),
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
	b.flash = vigor.NewFlashEffect(b.back, 0.15, ease.Linear, ease.Linear)
}

type Paddle struct {
	vigor.Image
	x int
}

func NewPaddle(x int, flip bool) *Paddle {
	p := &Paddle{
		x:     x,
		Image: *vigor.NewImage("paddle"),
	}
	p.PlaceOffScreen()
	if flip {
		p.FlipX()
	}
	return p
}

func (p *Paddle) PlaceOffScreen() {
	p.SetPos(float32(p.x), float32(screenHeight+p.Dim().Y))
}

func (p *Paddle) PlaceRandomly() {
	y := float32(rand.Intn(paddleMaxY) + paddleMinY)
	x := p.Pos().X
	p.TweenTo(x, y, 0.1, ease.Linear)
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

func (d *Dove) Die() {
	d.Show(false)
	d.SetMotion(false)
	d.SetPos(screenWidth/2, screenHeight/2)
}

func (d *Dove) Live() {
	d.Show(true)
	d.SetMotion(true)
	d.SetPos(screenWidth/2, screenHeight/2)
	d.SetVel(0, 0)
	d.SetAccel(0, 0)
	d.Scale(1, 1)
}

func (d *Dove) Init() {
	d.Live()
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
