package domain

import "strings"

// defaultSizePreset is the size used when the requested preset is missing or unsupported.
const defaultSizePreset = "md"

// sizePresets maps supported size names to badge dimensions in pixels.
// The map is read-only after initialization.
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
	preset = strings.ToLower(strings.TrimSpace(preset))

	if d, ok := sizePresets[preset]; ok {
		return d
	}
	return sizePresets[defaultSizePreset]
}
