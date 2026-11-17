package domain

import "strings"

// defaultTheme is the theme used when the requested theme is missing or unsupported.
const defaultTheme = "dark"

// themes is the set of supported theme names.
// It is also used to strip theme suffixes from badge asset names.
// The map is read-only after initialization.
var themes = map[string]struct{}{
	"dark":  {},
	"light": {},
}

func ValidateTheme(theme string) string {
	theme = strings.ToLower(strings.TrimSpace(theme))

	if _, ok := themes[theme]; ok {
		return theme
	}
	return defaultTheme
}

func TrimThemeSuffix(name string) string {
	for t := range themes {
		if before, found := strings.CutSuffix(name, "-"+t); found {
			return before
		}
	}
	return name
}
