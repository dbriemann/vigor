package vigor

import (
	"encoding/json"
	"io/ioutil"
)

// TODO: use better constraint than 'any'
func loadData[T any](fpath string) (T, error) {
	var t T

	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal(raw, &t); err != nil {
		return t, err
	}

	return t, nil
}

type AnimationSetConfig struct {
	Sections         []SectionConfig            `json:"sections"`
	Animations       map[string]AnimationConfig `json:"animations"`
	DefaultAnimation string                     `json:"defaultAnimation"`
}

type AnimationConfig struct {
	ImagePath  string `json:"imagePath"`
	SectionID  int    `json:"sectionId"`
	Frames     []int  `json:"frames"`
	DurationMS int    `json:"durationMs"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	EaseFunc   string `json:"easeFunc"`
	Loops      int    `json:"loops"`
}

type SectionConfig struct {
	Left    int `json:"left"`
	Top     int `json:"top"`
	Width   int `json:"width"`
	Height  int `json:"height"`
	Padding int `json:"padding"`
}
