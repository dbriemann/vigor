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
	ErrColumnMismatch   = fmt.Errorf("columns do not match with frame width and padding")
	ErrRowMismatch      = fmt.Errorf("rows do not match with frame height and padding")
	ErrFrameCountZero   = fmt.Errorf("the calculated frame count is zero")
	ErrUnknownEaseFunc  = fmt.Errorf("ease function name unknown")
	ErrNoAnimations     = fmt.Errorf("no animations defined")
	ErrNoSections       = fmt.Errorf("no sections defined")
	ErrUnknownAnimation = fmt.Errorf("unknown animation")
	ErrFileNotFound     = fmt.Errorf("file not found")
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

type AnimationSet struct {
	Sections        []Section
	Animations      map[string]*Animation
	ActiveAnimation *Animation
}

type ImageLoaderFunc func(path string) (*ebiten.Image, error)

func LoadAnimationSet(loaderFunc ImageLoaderFunc, cfgFile string) (*AnimationSet, error) {
	cfg, err := loadData[AnimationSetConfig](cfgFile)
	if err != nil {
		return nil, err
	}
	if len(cfg.Animations) == 0 {
		return nil, ErrNoAnimations
	}
	if len(cfg.Sections) == 0 {
		return nil, ErrNoSections
	}

	as := &AnimationSet{
		Sections:   []Section{},
		Animations: map[string]*Animation{},
	}
	for _, secCfg := range cfg.Sections {
		as.Sections = append(as.Sections, NewSection(secCfg.Left, secCfg.Top, secCfg.Width, secCfg.Height, secCfg.Padding))
	}
	for name, animCfg := range cfg.Animations {
		img, err := loaderFunc(animCfg.ImagePath)
		if err != nil {
			return nil, ErrFileNotFound
		}

		f, ok := easeFuncMappings[animCfg.EaseFunc]
		if !ok {
			return nil, ErrUnknownEaseFunc
		}
		a, err := NewAnimation(
			img,
			as.Sections[animCfg.SectionID],
			animCfg.Width,
			animCfg.Height,
			animCfg.Frames,
			time.Duration(animCfg.DurationMS*int(time.Millisecond)),
			animCfg.Loops,
			f,
		)
		if err != nil {
			return nil, err
		}
		as.Animations[name] = a
	}

	if err := as.Switch(cfg.DefaultAnimation); err != nil {
		return nil, err
	}

	return as, nil
}

// Switch changes the active animation and returns an ErrUnknownAnimation if the name does not exist.
func (as *AnimationSet) Switch(name string) error {
	if as.ActiveAnimation != nil {
		as.ActiveAnimation.Stop()
		// TODO reset?
	}

	anim, ok := as.Animations[name]
	if !ok {
		return ErrUnknownAnimation
	}
	as.ActiveAnimation = anim

	// TODO: start?

	return nil
}

func (as *AnimationSet) Run() {
	if as.ActiveAnimation != nil {
		as.ActiveAnimation.Run()
	}
}

func (as *AnimationSet) Stop() {
	if as.ActiveAnimation != nil {
		as.ActiveAnimation.Stop()
	}
}

func (as *AnimationSet) Reset() {
	if as.ActiveAnimation != nil {
		as.ActiveAnimation.Reset()
	}
}

func (as *AnimationSet) Update(dt time.Duration) {
	if as.ActiveAnimation != nil {
		as.ActiveAnimation.Update(dt)
	}
}

func (as *AnimationSet) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	if as.ActiveAnimation != nil {
		as.ActiveAnimation.Draw(target, op)
	}
}

type Animation struct {
	spriteSheet *ebiten.Image
	section     Section
	frameWidth  int
	frameHeight int
	sprites     []*ebiten.Image

	frames       []int
	duration     time.Duration
	tween        *gween.Tween
	currentFrame int
	lastFrame    int
	loops        int
	paused       bool
}

func NewAnimation(sheet *ebiten.Image, section Section, w, h int, frames []int, duration time.Duration, loops int, tweenFunc ease.TweenFunc) (*Animation, error) {
	a := &Animation{
		spriteSheet:  sheet,
		section:      section,
		frameWidth:   w,
		frameHeight:  h,
		frames:       frames,
		duration:     duration,
		currentFrame: 0,
		lastFrame:    0,
		paused:       true,
		loops:        loops,
	}

	// Calculate frame positions relative to the sprite sheet.
	a.sprites = []*ebiten.Image{}

	// NOTE: A padding larger than the frame will break this function.

	// Check if section and frame size are a fit.
	if a.spriteSheet.Bounds().Max.X%(a.frameWidth+a.section.padding) != a.section.padding {
		return nil, ErrColumnMismatch
	}
	if a.spriteSheet.Bounds().Max.Y%(a.frameHeight+a.section.padding) != a.section.padding {
		return nil, ErrRowMismatch
	}
	columns := a.spriteSheet.Bounds().Max.X / (a.frameWidth + a.section.padding)
	rows := a.spriteSheet.Bounds().Max.Y / (a.frameHeight + a.section.padding)

	numFrames := columns * rows
	if numFrames == 0 {
		return nil, ErrFrameCountZero
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			upperLeft := image.Point{
				X: a.section.left + (x+1)*a.section.padding + x*a.frameWidth,
				Y: a.section.top + (y+1)*a.section.padding + y*a.frameHeight,
			}
			subImg := a.spriteSheet.SubImage(image.Rect(
				upperLeft.X,
				upperLeft.Y,
				upperLeft.X+a.frameWidth,
				upperLeft.Y+a.frameHeight,
			)).(*ebiten.Image)
			a.sprites = append(a.sprites, subImg)
		}
	}

	a.tween = gween.New(0, float32(len(a.frames)-1), float32(a.duration.Seconds()), tweenFunc)

	return a, nil
}

// Run starts or resumes an animation.
func (a *Animation) Run() {
	a.paused = false
}

// Reset set the animation to frame zero, keeps running/paused state as is.
func (a *Animation) Reset() {
	a.currentFrame = 0
	a.lastFrame = 0
	a.tween.Reset()
}

// Stop pauses an animation at the current frame.
func (a *Animation) Stop() {
	a.paused = true
}

// Update selects the current frame to draw by applying delta time and the
// easing function.
func (a *Animation) Update(dt time.Duration) {
	if a.paused || len(a.frames) <= 1 {
		return
	}

	if a.loops == 0 {
		return
	}

	interpolation, finished := a.tween.Update(float32(dt.Seconds()))
	a.lastFrame = a.currentFrame
	frameIndex := int(math.Round(float64(interpolation)))
	a.currentFrame = a.frames[frameIndex]
	if finished && a.loops != 0 {
		a.tween.Reset()
		if a.loops > 0 {
			a.loops--
		}
	}
}

func (a *Animation) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	target.DrawImage(a.sprites[a.currentFrame], op)
}

func (a *Animation) SetFrames(frames []int) {
	a.frames = frames
}

func (a *Animation) SetDuration(dur time.Duration) {
	a.duration = dur
}

func (a *Animation) SetTweenFunc(f ease.TweenFunc) {
	a.currentFrame = 0
	a.lastFrame = 0
	a.tween = gween.New(0, float32(len(a.frames)-1), float32(a.duration.Seconds()), f)
}
