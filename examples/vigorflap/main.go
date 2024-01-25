package main

// TODO: attribution

import (
	_ "embed"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"strconv"

	"github.com/dbriemann/vigor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

const (
	screenWidth  = 160
	screenHeight = 240
)

var (
	dt                 = 1.0 / float32(ebiten.TPS())
	score      int     = 0
	highscore  int     = 0
	sideTop    float32 = 0
	sideBottom float32 = 0
)

type XY struct {
	X float32
	Y float32
}

type Dove struct {
	bbox       vigor.Rect[float32]
	vel        XY
	accel      XY
	flip       bool
	activeAnim *vigor.Animation
	animSail   *vigor.Animation
	animFlap   *vigor.Animation
}

func (d *Dove) Update() {
	d.bbox.Point.X += d.vel.X * dt
	d.bbox.Point.Y += d.vel.Y * dt

	d.vel.X += d.accel.X * dt
	d.vel.Y += d.accel.Y * dt

	if d.activeAnim == d.animFlap && d.activeAnim.Finished {
		d.activeAnim = d.animSail
		d.activeAnim.Reset()
		d.activeAnim.Run()
	}

	d.activeAnim.Update(dt)
}

func (d *Dove) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	d.activeAnim.Draw(target, op)
}

type Paddle struct {
	bbox  vigor.Rect[float32]
	tween *gween.Tween
	image *ebiten.Image
}

func (p *Paddle) PlaceRandom() {
	r := rand.Intn(int(sideBottom-sideTop-p.bbox.Dim.H)) + int(sideTop)
	p.tween = gween.New(p.bbox.Point.Y, float32(r), 0.2, ease.Linear)
}

func (p *Paddle) Update() {
	if p.tween == nil {
		return
	}
	y, done := p.tween.Update(dt)
	if done {
		p.tween = nil
		return
	}
	p.bbox.Point.Y = y
}

func (p *Paddle) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	op.GeoM.Translate(float64(p.bbox.Point.X), float64(p.bbox.Point.Y))
	target.DrawImage(p.image, op)
}

type Bouncer struct {
	bbox vigor.Rect[float32]
}

func (b *Bouncer) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	vector.DrawFilledRect(
		target,
		float32(b.bbox.Point.X),
		float32(b.bbox.Point.Y),
		float32(b.bbox.Dim.W),
		float32(b.bbox.Dim.H),
		color.RGBA{53, 53, 61, 240},
		false,
	)
	vector.DrawFilledRect(
		target,
		float32(b.bbox.Point.X+1),
		float32(b.bbox.Point.Y+1),
		float32(b.bbox.Dim.W-2),
		float32(b.bbox.Dim.H-2),
		color.RGBA{100, 106, 125, 225},
		false,
	)
}

type Game struct {
	man          vigor.ResourceManager
	background   *ebiten.Image
	dove         Dove
	spikes       *ebiten.Image
	paddle       *ebiten.Image
	paddleLeft   Paddle
	paddleRight  Paddle
	bouncerLeft  Bouncer
	bouncerRight Bouncer
	// TODO: sides/borders
}

func (g *Game) Init() {
	score = 0

	off := screenHeight + g.paddleLeft.bbox.Dim.H
	// TODO: why is setting Y not enough?
	g.paddleLeft.bbox.Point.Y = float32(off)
	g.paddleLeft.tween = nil
	g.paddleRight.bbox.Point.Y = float32(off)
	g.paddleRight.tween = nil

	g.dove.bbox.Point.X = screenWidth / 2
	g.dove.bbox.Point.Y = screenHeight / 2
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

	if g.dove.bbox.Intersects(g.bouncerLeft.bbox) {
		g.paddleRight.PlaceRandom()
		score++
		g.dove.vel.X *= -1
		g.dove.flip = false
	}

	if g.dove.bbox.Intersects(g.bouncerRight.bbox) {
		g.paddleLeft.PlaceRandom()
		score++
		g.dove.vel.X *= -1
		g.dove.flip = true
	}

	g.paddleLeft.Update()
	g.paddleRight.Update()
	g.dove.Update()

	if g.dove.bbox.Point.Y <= float32(g.spikes.Bounds().Dy()) ||
		g.dove.bbox.Point.Y+g.dove.bbox.Dim.H >= float32(screenHeight-g.spikes.Bounds().Dy()) ||
		g.dove.bbox.Intersects(g.paddleLeft.bbox) ||
		g.dove.bbox.Intersects(g.paddleRight.bbox) {
		g.Over()
	}

	return nil
}

func (g *Game) Over() {
	if score > highscore {
		highscore = score
	}
	score = 0
	// ly := g.paddleLeft.bbox.Point.Y
	// ry := g.paddleRight.bbox.Point.Y
	g.Init()
	// g.paddleLeft.tween = gween.New(ly, g.paddleLeft.bbox.Point.Y, 0.1, ease.Linear)
	// g.paddleRight.tween = gween.New(ry, g.paddleRight.bbox.Point.Y, 0.1, ease.Linear)
}

func (g *Game) Draw(target *ebiten.Image) {
	target.DrawImage(g.background, nil)

	ebitenutil.DebugPrintAt(target, strconv.Itoa(score), screenWidth/2-5, screenHeight-50)
	if highscore > 0 {
		ebitenutil.DebugPrintAt(target, "high: "+strconv.Itoa(highscore), screenWidth/2-20, 50)
	}

	topSpikesOp := &ebiten.DrawImageOptions{}
	topSpikesOp.GeoM.Scale(1, -1)
	topSpikesOp.GeoM.Translate(0, float64(g.spikes.Bounds().Dy()))
	target.DrawImage(g.spikes, topSpikesOp)

	bottomSpikesOp := &ebiten.DrawImageOptions{}
	bottomSpikesOp.GeoM.Translate(0, float64(screenHeight-g.spikes.Bounds().Dy()))
	target.DrawImage(g.spikes, bottomSpikesOp)

	g.bouncerLeft.Draw(target, nil)
	g.bouncerRight.Draw(target, nil)

	lpadOp := &ebiten.DrawImageOptions{}
	g.paddleLeft.Draw(target, lpadOp)

	rpadOp := &ebiten.DrawImageOptions{}
	rpadOp.GeoM.Scale(-1, 1)
	rpadOp.GeoM.Translate(float64(g.paddleRight.bbox.Dim.W), 0)
	g.paddleRight.Draw(target, rpadOp)

	doveOp := &ebiten.DrawImageOptions{}
	if g.dove.flip {
		doveOp.GeoM.Scale(-1, 1)
		doveOp.GeoM.Translate(float64(g.dove.bbox.Dim.W), 0)
	}
	doveOp.GeoM.Translate(float64(g.dove.bbox.Point.X), float64(g.dove.bbox.Point.Y))
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

	g.background = g.man.GetImageOrPanic("background")
	g.spikes = g.man.GetImageOrPanic("spikes")
	g.paddle = g.man.GetImageOrPanic("paddle")

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

	g.paddleLeft.image = g.paddle
	g.paddleRight.image = g.paddle
	g.paddleLeft.bbox.Dim.W = float32(g.paddle.Bounds().Dx())
	g.paddleLeft.bbox.Dim.H = float32(g.paddle.Bounds().Dy())
	g.paddleRight.bbox.Dim.W = float32(g.paddle.Bounds().Dx())
	g.paddleRight.bbox.Dim.H = float32(g.paddle.Bounds().Dy())

	g.Init()
	g.dove.bbox.Dim.W = float32(g.dove.activeAnim.FrameWidth)
	g.dove.bbox.Dim.H = float32(g.dove.activeAnim.FrameHeight)

	sideTop = float32(g.spikes.Bounds().Dy() + 2)
	sideBottom = float32(screenHeight - 2*g.spikes.Bounds().Dy() - 4)

	g.bouncerLeft.bbox.Dim.W = 4
	g.bouncerLeft.bbox.Point.X = 1
	g.bouncerLeft.bbox.Point.Y = sideTop
	g.bouncerLeft.bbox.Dim.H = sideBottom

	g.bouncerRight.bbox.Dim.W = 4
	g.bouncerRight.bbox.Point.X = screenWidth - g.bouncerRight.bbox.Dim.W - 1
	g.bouncerRight.bbox.Point.Y = sideTop
	g.bouncerRight.bbox.Dim.H = sideBottom

	g.paddleLeft.bbox.Point.X = float32(g.bouncerLeft.bbox.Dim.W + 2)
	g.paddleRight.bbox.Point.X = screenWidth - g.paddleRight.bbox.Dim.W - float32(g.bouncerRight.bbox.Dim.W) - 2

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
