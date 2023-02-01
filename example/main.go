package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
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
)

func GetFunctionName(i interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

var (
	// The awesome sprite is taken from here:
	// https://aamatniekss.itch.io/fantasy-knight-free-pixelart-animated-character
	// Go checkout the artist's pixel work.
	//
	//go:embed _AttackComboNoMovement.png
	imageFile []byte

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
	anim      *vigor.Animation
	funcIndex int
	millis    int
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.millis += 100
		g.recreate()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if g.millis >= 200 {
			g.millis -= 100
			g.recreate()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if g.funcIndex > 1 {
			g.funcIndex--
			g.recreate()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.funcIndex < len(easeFuncs)-1 {
			g.funcIndex++
			g.recreate()
		}
	}

	g.anim.Update(time.Second / time.Duration(ebiten.TPS())) // 1/60 sec

	return nil
}

func (g *Game) Draw(target *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.GeoM.Translate(-frameWidth/2, -frameHeight/2)
	g.anim.Draw(target, op)

	msg := fmt.Sprintf("Ease func: %s\nDuration: %d ms",
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
		millis:    500,
	}

	img, _, err := image.Decode(bytes.NewReader(imageFile))
	if err != nil {
		log.Fatal(err)
	}
	g.sheet = ebiten.NewImageFromImage(img)

	g.recreate()

	return g
}

func (g *Game) recreate() {
	section := vigor.NewSection(0, 0, frameWidth*10, frameHeight, 0)
	var err error
	g.anim, err = vigor.NewAnimation(
		g.sheet,
		section,
		frameWidth,
		frameHeight,
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		time.Duration(g.millis)*time.Millisecond,
		-1,
		easeFuncs[g.funcIndex],
	)
	if err != nil {
		panic(err)
	}
	g.anim.Start()
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
