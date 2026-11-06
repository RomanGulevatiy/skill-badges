package domain

import "strings"

const defaultSizePreset = "md"

var sizePresets = map[string]Dimensions{
	"sm": {Size: 32, Gap: 6},
	"md": {Size: 48, Gap: 8},
	"lg": {Size: 64, Gap: 12},
	"xl": {Size: 80, Gap: 16},
}

type Dimensions struct {
	Size int
	Gap  int
}

func GetDimensions(preset string) Dimensions {
	cleanPreset := strings.ToLower(strings.TrimSpace(preset))

	if d, ok := sizePresets[cleanPreset]; ok {
		return d
	}
	return sizePresets[defaultSizePreset]
}
