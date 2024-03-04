package main

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/dbriemann/vigor"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/tanema/gween/ease"
)

const (
	screenWidth    = 320
	screenHeight   = 240
	frameWidth     = 120
	frameHeight    = 80
	bgKnightsCount = 10

	ActionSlower input.Action = iota
	ActionFaster
	ActionPrevEaseFunc
	ActionNextEaseFunc
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

	keymap = input.Keymap{
		ActionSlower:       {input.KeyGamepadDown, input.KeyDown},
		ActionFaster:       {input.KeyGamepadUp, input.KeyUp},
		ActionPrevEaseFunc: {input.KeyGamepadLeft, input.KeyLeft},
		ActionNextEaseFunc: {input.KeyGamepadRight, input.KeyRight},
	}
)

type Game struct {
	bgKnights []*Knight
	knight    *Knight
	input     *input.Handler
	dur       time.Duration
	funcIndex int
}

func (g *Game) Init() {
	g.input = vigor.NewInputHandler(0, keymap)

	g.knight = NewKnight(0, 0)
	g.knight.SetPos(screenWidth/2-float32(frameWidth/2), screenHeight/2-float32(frameHeight/2))
	vigor.G.Add(g.knight)

	g.bgKnights = make([]*Knight, bgKnightsCount)
	for i := range bgKnightsCount {
		g.bgKnights[i] = NewKnight(0, 0)
		g.bgKnights[i].SetDuration(time.Millisecond * time.Duration(500+i*100))
		g.bgKnights[i].SetPos(float32(screenWidth*i)/bgKnightsCount, float32(screenHeight/2+frameHeight))
		g.bgKnights[i].Scale(0.3, 0.3)
		vigor.G.Add(g.bgKnights[i])
	}
}

func (g *Game) Update() {
	if g.input.ActionIsJustPressed(ActionFaster) {
		g.dur += 100 * time.Millisecond
		g.knight.SetDuration(g.dur)
	} else if g.input.ActionIsJustPressed(ActionSlower) {
		if g.dur >= 200*time.Millisecond {
			g.dur -= 100 * time.Millisecond
			g.knight.SetDuration(g.dur)
		}
	} else if g.input.ActionIsJustPressed(ActionPrevEaseFunc) {
		if g.funcIndex > 0 {
			g.funcIndex--
			g.knight.SetTweenFunc(easeFuncs[g.funcIndex])
		}
	} else if g.input.ActionIsJustPressed(ActionNextEaseFunc) {
		if g.funcIndex < len(easeFuncs)-1 {
			g.funcIndex++
			g.knight.SetTweenFunc(easeFuncs[g.funcIndex])
		}
	}
	// fmt.Println("ease func:", GetFunctionName(easeFuncs[g.funcIndex]))
	// fmt.Println("duration:", g.dur)
	vigor.DebugPrintf("Ease func: %s (left/right arrows)\nDuration: %s (up/down arrows)",
		GetFunctionName(easeFuncs[g.funcIndex]),
		g.dur,
	)
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

type Knight struct {
	vigor.Sprite
}

func (k *Knight) Update() {
	k.Sprite.Update()
}

func NewKnight(x, y int) *Knight {
	k := &Knight{
		Sprite: *vigor.NewSprite("knight_attack1"),
	}
	k.SetPos(float32(x), float32(y))
	return k
}

func main() {
	vigor.SetWindowSize(4*screenWidth, 4*screenHeight)
	g := Game{
		dur:       700 * time.Millisecond,
		funcIndex: 0,
	}

	vigor.InitGame(&g)

	vigor.RunGame()
}
