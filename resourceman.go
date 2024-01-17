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

	Sections           map[string]Section
	AnimationTemplates map[string]*AnimationTemplate
}

func NewResourceManager() ResourceManager {
	r := ResourceManager{
		Images:             map[string]*ebiten.Image{},
		Sections:           map[string]Section{},
		AnimationTemplates: map[string]*AnimationTemplate{},
	}
	return r
}

func (r *ResourceManager) LoadConfig(fname string) error {
	cfg, err := loadConfigData[ResourceConfig](fname)
	if err != nil {
		return err
	}

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
	for animName, template := range cfg.Animations {
		// Check if image exists in loaded assets.
		imgName := template.ImageName
		img, ok := r.Images[imgName]
		if !ok {
			return fmt.Errorf("%w: %s", ErrImageNotLoaded, imgName)
		}

		f, ok := easeFuncMappings[template.EaseFunc]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownEaseFunc, template.EaseFunc)
		}
		a, err := NewAnimationTemplate(
			img,
			r.Sections[template.SectionName],
			template.Width,
			template.Height,
			template.Frames,
			time.Duration(template.DurationMS*int(time.Millisecond)),
			template.Loops,
			f,
		)
		if err != nil {
			return err
		}
		r.AnimationTemplates[animName] = a
	}

	return nil
}
