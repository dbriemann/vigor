package main

// The awesome knight sprite is taken from here:
// https://aamatniekss.itch.io/fantasy-knight-free-pixelart-animated-character
// Go checkout the artist's pixel work.

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/dbriemann/vigor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tanema/gween/ease"
)

const (
	screenWidth  = 320
	screenHeight = 240
	frameWidth   = 120
	frameHeight  = 80

	bgKnightsCount = 10
)

func GetFunctionName(i interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

var (
	easeFuncs = []ease.TweenFunc{
		ease.Linear,
		ease.InQuad, ease.OutQuad, ease.InOutQuad, ease.OutInQuad,
		ease.InCubic, ease.OutCubic, ease.InOutCubic, ease.OutInCubic,
		ease.InQuart, ease.OutQuart, ease.InOutQuart, ease.OutInQuart,
		ease.InQuint, ease.OutQuint, ease.InOutQuint, ease.OutInQuint,
		ease.InExpo, ease.OutExpo, ease.InOutExpo, ease.OutInExpo,
		ease.InSine, ease.OutSine, ease.InOutSine, ease.OutInSine,
		ease.InCirc, ease.OutCirc, ease.InOutCirc, ease.OutInCirc,
		ease.InBack, ease.OutBack, ease.InOutBack, ease.OutInBack,
		ease.InBounce, ease.OutBounce, ease.InOutBounce, ease.OutInBounce,
		ease.InElastic, ease.OutElastic, ease.InOutElastic, ease.OutInElastic,
	}

	dt = 1.0 / float32(ebiten.TPS())
)

type Game struct {
	man        vigor.AssetManager
	knightAnim *vigor.Animation
	bgKnights  []*vigor.Animation
	funcIndex  int
	dur        time.Duration
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.dur += 100 * time.Millisecond
		g.applySettings()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if g.dur >= 200*time.Millisecond {
			g.dur -= 100 * time.Millisecond
			g.applySettings()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if g.funcIndex > 0 {
			g.funcIndex--
			g.applySettings()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.funcIndex < len(easeFuncs)-1 {
			g.funcIndex++
			g.applySettings()
		}
	}

	g.knightAnim.Update(dt)
	for i := 0; i < bgKnightsCount; i++ {
		g.bgKnights[i].Update(dt)
	}

	return nil
}

func (g *Game) Draw(target *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(-frameWidth/2, -frameHeight/2)
	g.knightAnim.Draw(target, op)

	for i := 0; i < bgKnightsCount; i++ {
		op.GeoM.Reset()
		op.GeoM.Scale(0.3, 0.3)
		op.GeoM.Translate(float64(screenWidth*i/bgKnightsCount), screenHeight/2+frameHeight)
		g.bgKnights[i].Draw(target, op)
	}

	msg := fmt.Sprintf("Ease func: %s (left/right arrows)\nDuration: %s (up/down arrows)",
		GetFunctionName(easeFuncs[g.funcIndex]),
		g.dur,
	)
	ebitenutil.DebugPrint(target, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{
		funcIndex: 0,
		man:       vigor.NewAssetManager(),
		bgKnights: make([]*vigor.Animation, bgKnightsCount),
	}

	if err := g.man.LoadConfig("assets/config.json"); err != nil {
		panic(err)
	}

	a, err := vigor.NewAnimation(g.man.AnimationTemplates["knight_attack1"])
	if err != nil {
		panic(err)
	}
	g.knightAnim = a

	g.dur = g.knightAnim.Duration

	g.knightAnim.Run()

	for i := 0; i < bgKnightsCount; i++ {
		a, err := vigor.NewAnimation(g.man.AnimationTemplates["knight_attack1"])
		if err != nil {
			panic(err)
		}
		g.bgKnights[i] = a
		g.bgKnights[i].SetDuration(time.Millisecond * time.Duration(500+i*100))
		g.bgKnights[i].Run()
	}

	return g
}

func (g *Game) applySettings() {
	a := g.knightAnim
	a.SetDuration(g.dur)
	a.SetTweenFunc(easeFuncs[g.funcIndex])
	a.InitTween()
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
