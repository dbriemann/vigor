package vigor

import (
	"encoding/json"
	"os"
)

func loadConfigData[T any](fpath string) (T, error) {
	var t T

	raw, err := os.ReadFile(fpath)
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal(raw, &t); err != nil {
		return t, err
	}

	return t, nil
}

type ResourceConfig struct {
	// TODO: audio
	// TODO: others
	Images       map[string]string          `json:"images"`
	Sections     map[string]SectionConfig   `json:"sections"`
	Animations   map[string]AnimationConfig `json:"animations"`
	ResourceRoot string                     `json:"resourceRoot"`
}

type AnimationConfig struct {
	ImageName   string  `json:"imageName"`
	SectionName string  `json:"sectionName"`
	EaseFunc    string  `json:"easeFunc"`
	Frames      []int   `json:"frames"`
	Duration    float64 `json:"duration"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Looped      bool    `json:"looped"`
}

type SectionConfig struct {
	Left    int `json:"left"`
	Top     int `json:"top"`
	Width   int `json:"width"`
	Height  int `json:"height"`
	Padding int `json:"padding"`
}
