package domain

import "strings"

var aliases = map[string]string{
	"go":       "golang",
	"js":       "javascript",
	"ts":       "typescript",
	"alpine":   "alpinelinux",
	"arch":     "archlinux",
	"idea":     "intellijidea",
	"postgres": "postgresql",
}

func ResolveAlias(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	if canonical, ok := aliases[name]; ok {
		return canonical
	}
	return name
}
