package domain

import "strings"

const defaultTheme = "dark"

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
