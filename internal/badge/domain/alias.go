package domain

import "strings"

var Aliases = map[string]string{
	"go":       "golang",
	"js":       "javascript",
	"ts":       "typescript",
	"alpine":   "alpinelinux",
	"arch":     "archlinux",
	"idea":     "intellijidea",
	"postgres": "postgresql",
}

func ResolveAlias(name string) string {
	lowerName := strings.ToLower(name)

	if canonical, ok := Aliases[lowerName]; ok {
		return canonical
	}
	return lowerName
}
