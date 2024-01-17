package vigor

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

// TODO: add callbacks to update, draw, etc.

var (
	ErrColumnMismatch     = fmt.Errorf("columns do not match with frame width and padding")
	ErrRowMismatch        = fmt.Errorf("rows do not match with frame height and padding")
	ErrFrameCountZero     = fmt.Errorf("the calculated frame count is zero")
	ErrFrameExceedsBounds = fmt.Errorf("frame index exceeds section bounds")
	ErrUnknownEaseFunc    = fmt.Errorf("ease function name unknown")
	ErrNoAnimations       = fmt.Errorf("no animations defined")
	ErrNoSections         = fmt.Errorf("no sections defined")
	ErrUnknownAnimation   = fmt.Errorf("unknown animation")
	ErrFileNotFound       = fmt.Errorf("file not found")
	ErrImageNotLoaded     = fmt.Errorf("image not loaded")
)

type Section struct {
	left    int
	top     int
	width   int
	height  int
	padding int
}

func NewSection(left, top, width, height, padding int) Section {
	s := Section{
		left:    left,
		top:     top,
		width:   width,
		height:  height,
		padding: padding,
	}
	return s
}

func (s *Section) Bounds() image.Rectangle {
	return image.Rect(s.left, s.top, s.left+s.width, s.top+s.height)
}

type AnimationTemplate struct {
	Section     Section
	FrameWidth  int
	FrameHeight int
	Sheet       *ebiten.Image
	Sprites     []*ebiten.Image
	Frames      []int
	Duration    time.Duration
	EaseFunc    ease.TweenFunc
	Loops       int
}

func NewAnimationTemplate(sheet *ebiten.Image, section Section, w, h int, frames []int, duration time.Duration, loops int, easeFunc ease.TweenFunc) (*AnimationTemplate, error) {
	t := AnimationTemplate{
		Section:     section,
		FrameWidth:  w,
		FrameHeight: h,
		Frames:      frames,
		Sheet:       sheet,
		Sprites:     []*ebiten.Image{},
		Duration:    duration,
		EaseFunc:    easeFunc,
		Loops:       loops,
	}

	// Calculate frame positions relative to the sprite sheet.
	// NOTE: A padding larger than the frame will break this function.

	// Check if section and frame size are a fit.
	if sheet.Bounds().Max.X%(t.FrameWidth+t.Section.padding) != t.Section.padding {
		return nil, ErrColumnMismatch
	}
	if sheet.Bounds().Max.Y%(t.FrameHeight+t.Section.padding) != t.Section.padding {
		return nil, ErrRowMismatch
	}
	columns := sheet.Bounds().Max.X / (t.FrameWidth + t.Section.padding)
	rows := sheet.Bounds().Max.Y / (t.FrameHeight + t.Section.padding)

	numFrames := columns * rows
	if numFrames == 0 {
		return nil, ErrFrameCountZero
	}
	for _, f := range t.Frames {
		if f >= numFrames {
			return nil, ErrFrameExceedsBounds
		}
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			upperLeft := image.Point{
				X: t.Section.left + (x+1)*t.Section.padding + x*t.FrameWidth,
				Y: t.Section.top + (y+1)*t.Section.padding + y*t.FrameHeight,
			}
			subImg := sheet.SubImage(image.Rect(
				upperLeft.X,
				upperLeft.Y,
				upperLeft.X+t.FrameWidth,
				upperLeft.Y+t.FrameHeight,
			)).(*ebiten.Image)
			t.Sprites = append(t.Sprites, subImg)
		}
	}

	return &t, nil
}

type Animation struct {
	*AnimationTemplate
	Tween        gween.Tween
	CurrentFrame int
	LastFrame    int
	Paused       bool
}

func NewAnimation(template *AnimationTemplate) (*Animation, error) {
	a := &Animation{
		CurrentFrame:      0,
		LastFrame:         0,
		Paused:            true,
		AnimationTemplate: template,
	}

	a.Tween = *gween.New(0, float32(len(a.Frames)-1), float32(a.Duration.Seconds()), a.EaseFunc)

	return a, nil
}

// Run starts or resumes an animation.
func (a *Animation) Run() {
	a.Paused = false
}

// Reset set the animation to frame zero, keeps running/paused state as is.
func (a *Animation) Reset() {
	a.CurrentFrame = 0
	a.LastFrame = 0
	a.Tween.Reset()
}

// Stop pauses an animation at the current frame.
func (a *Animation) Stop() {
	a.Paused = true
}

// Update selects the current frame to draw considering the easing function.
func (a *Animation) Update() {
	if a.Paused || len(a.Frames) <= 1 {
		return
	}

	if a.Loops == 0 {
		return
	}

	dt := 1.0 / float32(ebiten.TPS())
	interpolation, finished := a.Tween.Update(dt)
	a.LastFrame = a.CurrentFrame
	frameIndex := int(math.Round(float64(interpolation)))
	a.CurrentFrame = a.Frames[frameIndex]
	if finished && a.Loops != 0 {
		a.Tween.Reset()
		if a.Loops > 0 {
			a.Loops--
		}
	}
}

func (a *Animation) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	target.DrawImage(a.Sprites[a.CurrentFrame], op)
}

func (a *Animation) SetFrames(frames []int) {
	a.Frames = frames
}

func (a *Animation) SetDuration(dur time.Duration) {
	a.Duration = dur
}

func (a *Animation) SetTweenFunc(f ease.TweenFunc) {
	a.CurrentFrame = 0
	a.LastFrame = 0
	a.Tween = *gween.New(0, float32(len(a.Frames)-1), float32(a.Duration.Seconds()), f)
}
