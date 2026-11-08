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
	"maven":    "apachmaven",
	"jwt":      "jsonwebtoken",
	"tailwind": "tailwindcss",
	"vue":      "vuejs",
	"zen":      "zenbrowser",
	"cs":       "csharp",
	"css":      "css3",
	"html":     "html5",
	"rs":       "rust",
	"f":        "fortran",
	"kt":       "kotlin",
	"py":       "python",
	"rb":       "ruby",
}

func ResolveAlias(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	if canonical, ok := aliases[name]; ok {
		return canonical
	}
	return name
}
