package domain

import "strings"

const defaultTheme = "dark"

var themes = map[string]struct{}{
	"dark":  {},
	"light": {},
}

func ValidateTheme(theme string) string {
	cleanTheme := strings.ToLower(strings.TrimSpace(theme))

	if _, ok := themes[cleanTheme]; ok {
		return cleanTheme
	}
	return defaultTheme
}

func AvailableThemes() []string {
	list := make([]string, 0, len(themes))

	for t := range themes {
		list = append(list, t)
	}
	return list
}
