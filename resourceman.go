package vigor

import (
	"fmt"
	"image"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type ResourceManager struct {
	RootPath string
	Images   map[string]*ebiten.Image
	// TODO: audio
	// TODO: others

	Sections   map[string]Section
	Animations map[string]*Animation
}

func NewResourceManager() ResourceManager {
	r := ResourceManager{
		Images:     map[string]*ebiten.Image{},
		Sections:   map[string]Section{},
		Animations: map[string]*Animation{},
	}
	return r
}

func (r *ResourceManager) LoadConfig(fname string) error {
	cfg, err := loadData[ResourceConfig](fname)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", cfg)

	// TODO: use data paths from resource dir of project.
	// e.g. images/bla.png with resources at ROOT/resources will
	// translate to ROOT/resources/images/bla.png and use this as name.

	r.RootPath = cfg.ResourceRoot

	// Load images.
	for path, name := range cfg.Images {
		f, err := os.Open(r.RootPath + path) // TODO: join
		if err != nil {
			return err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		ebImg := ebiten.NewImageFromImage(img)

		r.Images[name] = ebImg
	}

	// TODO: audio
	// TODO: others

	// Init sections.
	for name, sec := range cfg.Sections {
		r.Sections[name] = NewSection(sec.Left, sec.Top, sec.Width, sec.Height, sec.Padding)
	}

	// Load animation templates.
	for animName, animation := range cfg.Animations {
		// Check if image exists in loaded assets.
		imgName := animation.ImageName
		img, ok := r.Images[imgName]
		if !ok {
			return fmt.Errorf("%w: %s", ErrImageNotLoaded, imgName)
		}

		f, ok := easeFuncMappings[animation.EaseFunc]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownEaseFunc, animation.EaseFunc)
		}
		a, err := NewAnimation(
			img,
			r.Sections[animation.SectionName],
			animation.Width,
			animation.Height,
			animation.Frames,
			time.Duration(animation.DurationMS*int(time.Millisecond)),
			animation.Loops,
			f,
		)
		if err != nil {
			return err
		}
		r.Animations[animName] = a
	}

	return nil
}
