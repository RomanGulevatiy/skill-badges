package domain

import "strings"

// aliases maps short user-facing badge names to their canonical asset names.
// The map is read-only after initialization.
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
	"c":        "clang",
	"cs":       "csharp",
	"css":      "css3",
	"html":     "html5",
	"rs":       "rust",
	"f":        "fortran",
	"kt":       "kotlin",
	"py":       "python",
	"rb":       "ruby",
	"as":       "assemblyscript",
	"shell":    "bash",
	"clj":      "clojure",
	"cr":       "crystal",
	"ex":       "elixir",
	"hs":       "haskell",
	"sc":       "scala",
	"wasm":     "webassembly",
}

func ResolveAlias(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	if canonical, ok := aliases[name]; ok {
		return canonical
	}
	return name
}
