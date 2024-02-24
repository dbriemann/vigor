package vigor

import (
	"fmt"
	"image"
	"os"
	"path"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: should this be private in the vigor package?
type AssetManager struct {
	// TODO: audio
	// TODO: others
	Images             map[string]*ebiten.Image
	Sections           map[string]Section
	AnimationTemplates map[string]*AnimationTemplate
	RootPath           string
}

func NewAssetManager() AssetManager {
	r := AssetManager{
		Images:             map[string]*ebiten.Image{},
		Sections:           map[string]Section{},
		AnimationTemplates: map[string]*AnimationTemplate{},
	}
	return r
}

func (r *AssetManager) LoadConfig(fname string) error {
	cfg, err := loadConfigData[ResourceConfig](fname)
	if err != nil {
		return err
	}

	r.RootPath = cfg.ResourceRoot

	for relPath, name := range cfg.Images {
		f, err := os.Open(path.Join(r.RootPath, relPath))
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

	for name, sec := range cfg.Sections {
		r.Sections[name] = NewSection(sec.Left, sec.Top, sec.Width, sec.Height, sec.Padding)
	}

	for animName, template := range cfg.Animations {
		imgName := template.ImageName
		img, ok := r.Images[imgName]
		if !ok {
			return fmt.Errorf("%w: %s", ErrImageNotLoaded, imgName)
		}

		if template.EaseFunc == "" {
			template.EaseFunc = "Linear"
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
			time.Duration(template.Duration*float64(time.Second)),
			template.Looped,
			f,
		)
		if err != nil {
			return err
		}
		r.AnimationTemplates[animName] = a
	}

	return nil
}

func (r *AssetManager) GetImageOrPanic(name string) *ebiten.Image {
	img, ok := r.Images[name]
	if !ok {
		panic(fmt.Sprintf("could not load image %s from asset manager: does not exist", name))
	}
	return img
}

func (r *AssetManager) GetAnimTemplateOrPanic(name string) *AnimationTemplate {
	templ, ok := r.AnimationTemplates[name]
	if !ok {
		panic(fmt.Sprintf("could not load animation template %s from asset manager: does not exist", name))
	}
	return templ
}
