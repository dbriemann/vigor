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
	ResourceRoot string            `json:"resourceRoot"`
	Images       map[string]string `json:"images"`
	// TODO: audio
	// TODO: others

	Sections   map[string]SectionConfig   `json:"sections"`
	Animations map[string]AnimationConfig `json:"animations"`
}

type AnimationConfig struct {
	ImageName   string `json:"imageName"`
	SectionName string `json:"sectionName"`
	Frames      []int  `json:"frames"`
	DurationMS  int    `json:"durationMs"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	EaseFunc    string `json:"easeFunc"`
	Loops       int    `json:"loops"`
}

type SectionConfig struct {
	Left    int `json:"left"`
	Top     int `json:"top"`
	Width   int `json:"width"`
	Height  int `json:"height"`
	Padding int `json:"padding"`
}
