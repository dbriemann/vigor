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

// TODO: add callbacks?
// after update, draw, loop?
// see what is needed first.

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
	ErrTemplateNotFound   = fmt.Errorf("template not found")
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
	Sheet       *ebiten.Image
	EaseFunc    ease.TweenFunc
	Sprites     []*ebiten.Image
	Frames      []int
	Section     Section
	FrameWidth  int
	FrameHeight int
	Duration    time.Duration
	Looped      bool
}

// TODO: simplify animations, remove sprite sheets and sections, just use images for each animation.

func NewAnimationTemplate(sheet *ebiten.Image, section Section, w, h int, frames []int, duration time.Duration, looped bool, easeFunc ease.TweenFunc) (*AnimationTemplate, error) {
	t := AnimationTemplate{
		Section:     section,
		FrameWidth:  w,
		FrameHeight: h,
		Frames:      frames,
		Sheet:       sheet,
		Sprites:     []*ebiten.Image{},
		Duration:    duration,
		EaseFunc:    easeFunc,
		Looped:      looped,
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
	Tween    *gween.Tween
	Frame    int
	Paused   bool
	Finished bool
}

func NewAnimation(template *AnimationTemplate) (*Animation, error) {
	if template == nil {
		return nil, ErrTemplateNotFound
	}
	if len(template.Frames) == 0 {
		return nil, ErrFrameCountZero
	}
	a := &Animation{
		AnimationTemplate: template,
		Frame:             template.Frames[0],
		Paused:            true,
		Finished:          len(template.Frames) == 1,
	}

	a.InitTween()

	return a, nil
}

// Run starts or resumes an animation.
func (a *Animation) Run() {
	a.Paused = false
}

// Reset set the animation to frame zero, keeps running/paused state as is.
func (a *Animation) Reset() {
	a.Frame = a.Frames[0]
	a.Tween.Reset()
}

// Stop pauses an animation at the current frame.
func (a *Animation) Stop() {
	a.Paused = true
}

// Update selects the current frame to draw considering the easing function.
func (a *Animation) Update(dt float32) {
	if a.Paused || a.Finished {
		return
	}

	interpolation, finished := a.Tween.Update(dt)
	frameIndex := int(math.Round(float64(interpolation)))
	a.Frame = a.Frames[frameIndex]

	if finished && a.Looped {
		a.Reset()
		return
	}
	a.Finished = finished
}

func (a *Animation) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	target.DrawImage(a.Sprites[a.Frame], op)
}

func (a *Animation) SetFrames(frames []int) {
	a.Frames = frames
	if len(a.Frames) == 1 {
		a.Finished = true
	}
}

func (a *Animation) SetDuration(duration time.Duration) {
	a.Duration = duration
}

func (a *Animation) SetTweenFunc(f ease.TweenFunc) {
	a.EaseFunc = f
}

func (a *Animation) InitTween() {
	a.Tween = gween.New(0, float32(len(a.Frames)-1), float32(a.Duration.Seconds()), a.EaseFunc)
}
