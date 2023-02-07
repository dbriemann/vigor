package main

import (
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
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
)

func GetFunctionName(i interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

var (
	// The awesome sprite is taken from here:
	// https://aamatniekss.itch.io/fantasy-knight-free-pixelart-animated-character
	// Go checkout the artist's pixel work.

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
)

type Game struct {
	sheet     *ebiten.Image
	animSet   *vigor.AnimationSet
	funcIndex int
	millis    int
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.millis += 100
		g.applySettings()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if g.millis >= 200 {
			g.millis -= 100
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

	g.animSet.Update(time.Second / time.Duration(ebiten.TPS())) // 1/60 sec

	return nil
}

func (g *Game) Draw(target *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(-frameWidth/2, -frameHeight/2)
	g.animSet.Draw(target, op)

	msg := fmt.Sprintf("Ease func: %s (left/right arrows)\nDuration: %d ms (up/down arrows)",
		GetFunctionName(easeFuncs[g.funcIndex]),
		g.millis,
	)
	ebitenutil.DebugPrint(target, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{
		funcIndex: 0,
		millis:    700,
	}

	loader := func(path string) (*ebiten.Image, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return nil, err
		}
		ebImg := ebiten.NewImageFromImage(img)
		return ebImg, nil
	}

	// Load set of animations from config file.
	aSet, err := vigor.LoadAnimationSet(loader, "anim_data.json")
	if err != nil {
		panic(err)
	}
	g.animSet = aSet

	g.animSet.Run()

	return g
}

func (g *Game) applySettings() {
	a := g.animSet.ActiveAnimation
	a.SetDuration(time.Duration(g.millis) * time.Millisecond)
	a.SetTweenFunc(easeFuncs[g.funcIndex])
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
