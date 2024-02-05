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
	"github.com/hajimehoshi/ebiten/v2/colorm"
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

type Dove struct {
	bbox       vigor.Rect[float32]
	vel        vigor.Vec2[float32]
	accel      vigor.Vec2[float32]
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
	r := rand.Intn(int(sideBottom-sideTop-p.bbox.Dim.Y)) + int(sideTop)
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
	imageBack  *ebiten.Image
	imageFront *ebiten.Image
	bbox       vigor.Rect[float32]
	flash      *vigor.FlashEffect
}

func (b *Bouncer) Draw(target *ebiten.Image) {
	opBack := &ebiten.DrawImageOptions{}
	opBack.GeoM.Translate(float64(b.bbox.Point.X), float64(b.bbox.Point.Y))
	opFront := &ebiten.DrawImageOptions{}
	opFront.GeoM.Translate(float64(b.bbox.Point.X+1), float64(b.bbox.Point.Y+1))

	target.DrawImage(b.imageBack, opBack)
	opFlash := &colorm.DrawImageOptions{}
	opFlash.GeoM.Translate(float64(b.bbox.Point.X), float64(b.bbox.Point.Y))
	b.flash.Draw(target, opFlash)
	target.DrawImage(b.imageFront, opFront)
}

func (b *Bouncer) Update() {
	b.flash.Update(dt)
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
	shake        *vigor.ShakeEffect
	deathAnim    bool
}

func (g *Game) Init() {
	score = 0

	off := screenHeight + g.paddleLeft.bbox.Dim.Y
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
	g.deathAnim = false
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
		g.bouncerLeft.flash.Reset()
		g.bouncerLeft.flash.Start()
		g.paddleRight.PlaceRandom()
		score++
		g.dove.vel.X *= -1
		g.dove.flip = false
	}

	if g.dove.bbox.Intersects(g.bouncerRight.bbox) {
		g.bouncerRight.flash.Reset()
		g.bouncerRight.flash.Start()
		g.paddleLeft.PlaceRandom()
		score++
		g.dove.vel.X *= -1
		g.dove.flip = true
	}

	g.paddleLeft.Update()
	g.paddleRight.Update()
	g.dove.Update()

	if g.dove.bbox.Point.Y <= float32(g.spikes.Bounds().Dy()) ||
		g.dove.bbox.Point.Y+g.dove.bbox.Dim.Y >= float32(screenHeight-g.spikes.Bounds().Dy()) ||
		g.dove.bbox.Intersects(g.paddleLeft.bbox) ||
		g.dove.bbox.Intersects(g.paddleRight.bbox) {
		g.Over()
	}

	g.bouncerLeft.Update()
	g.bouncerRight.Update()

	if g.deathAnim {
		g.dove.bbox.Point.X = 6666
		g.deathAnim = !g.shake.Update(dt)
		if !g.deathAnim {
			g.Init()
		}
	}

	return nil
}

func (g *Game) Over() {
	if g.deathAnim {
		return
	}

	if score > highscore {
		highscore = score
	}

	g.shake.Reset()
	g.shake.Start()
	g.deathAnim = true

	score = 0
}

func (g *Game) Draw(target *ebiten.Image) {
	scene := ebiten.NewImage(target.Bounds().Dx(), target.Bounds().Dy())
	scene.DrawImage(g.background, nil)

	ebitenutil.DebugPrintAt(scene, strconv.Itoa(score), screenWidth/2-5, screenHeight-50)
	if highscore > 0 {
		ebitenutil.DebugPrintAt(scene, "high: "+strconv.Itoa(highscore), screenWidth/2-20, 50)
	}

	topSpikesOp := &ebiten.DrawImageOptions{}
	topSpikesOp.GeoM.Scale(1, -1)
	topSpikesOp.GeoM.Translate(0, float64(g.spikes.Bounds().Dy()))
	scene.DrawImage(g.spikes, topSpikesOp)

	bottomSpikesOp := &ebiten.DrawImageOptions{}
	bottomSpikesOp.GeoM.Translate(0, float64(screenHeight-g.spikes.Bounds().Dy()))
	scene.DrawImage(g.spikes, bottomSpikesOp)

	g.bouncerLeft.Draw(scene)
	g.bouncerRight.Draw(scene)

	lpadOp := &ebiten.DrawImageOptions{}
	g.paddleLeft.Draw(scene, lpadOp)

	rpadOp := &ebiten.DrawImageOptions{}
	rpadOp.GeoM.Scale(-1, 1)
	rpadOp.GeoM.Translate(float64(g.paddleRight.bbox.Dim.X), 0)
	g.paddleRight.Draw(scene, rpadOp)

	doveOp := &ebiten.DrawImageOptions{}
	if g.dove.flip {
		doveOp.GeoM.Scale(-1, 1)
		doveOp.GeoM.Translate(float64(g.dove.bbox.Dim.X), 0)
	}
	doveOp.GeoM.Translate(float64(g.dove.bbox.Point.X), float64(g.dove.bbox.Point.Y))
	g.dove.Draw(scene, doveOp)

	sceneOp := &ebiten.DrawImageOptions{}
	g.shake.Apply(sceneOp)
	target.DrawImage(scene, sceneOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{
		man:       vigor.NewResourceManager(),
		shake:     vigor.NewShakeEffect(0.5, 8, 8),
		deathAnim: false,
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
	g.paddleLeft.bbox.Dim.X = float32(g.paddle.Bounds().Dx())
	g.paddleLeft.bbox.Dim.Y = float32(g.paddle.Bounds().Dy())
	g.paddleRight.bbox.Dim.X = float32(g.paddle.Bounds().Dx())
	g.paddleRight.bbox.Dim.Y = float32(g.paddle.Bounds().Dy())

	g.Init()
	g.dove.bbox.Dim.X = float32(g.dove.activeAnim.FrameWidth)
	g.dove.bbox.Dim.Y = float32(g.dove.activeAnim.FrameHeight)

	sideTop = float32(g.spikes.Bounds().Dy() + 2)
	sideBottom = float32(screenHeight - 2*g.spikes.Bounds().Dy() - 4)

	g.bouncerLeft.bbox.Dim.X = 4
	g.bouncerLeft.bbox.Point.X = 1
	g.bouncerLeft.bbox.Point.Y = sideTop
	g.bouncerLeft.bbox.Dim.Y = sideBottom

	g.bouncerRight.bbox.Dim.X = 4
	g.bouncerRight.bbox.Point.X = screenWidth - g.bouncerRight.bbox.Dim.X - 1
	g.bouncerRight.bbox.Point.Y = sideTop
	g.bouncerRight.bbox.Dim.Y = sideBottom

	g.paddleLeft.bbox.Point.X = float32(g.bouncerLeft.bbox.Dim.X + 2)
	g.paddleRight.bbox.Point.X = screenWidth - g.paddleRight.bbox.Dim.X - float32(g.bouncerRight.bbox.Dim.X) - 2

	// Create bouncer sprite
	bouncerBackImg := ebiten.NewImage(4, int(sideBottom))
	vector.DrawFilledRect(
		bouncerBackImg,
		float32(0),
		float32(0),
		float32(4),
		float32(sideBottom),
		color.RGBA{53, 53, 61, 240},
		false,
	)
	bouncerFrontImg := ebiten.NewImage(4-2, int(sideBottom)-2)
	vector.DrawFilledRect(
		bouncerFrontImg,
		float32(0),
		float32(0),
		float32(2),
		float32(sideBottom-2),
		color.RGBA{100, 106, 125, 225},
		false,
	)
	g.bouncerLeft.imageBack = bouncerBackImg
	g.bouncerLeft.imageFront = bouncerFrontImg
	g.bouncerRight.imageBack = bouncerBackImg
	g.bouncerRight.imageFront = bouncerFrontImg

	g.bouncerLeft.flash = vigor.NewFlashEffect(g.bouncerLeft.imageBack, 0.2, ease.Linear, ease.Linear)
	g.bouncerRight.flash = vigor.NewFlashEffect(g.bouncerRight.imageBack, 0.2, ease.Linear, ease.Linear)

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
