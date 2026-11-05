package domain

const DefaultTheme = "dark"

var themes = map[string]struct{}{
	"dark":  {},
	"light": {},
}

func IsValidTheme(theme string) (ok bool) {
	_, ok = themes[theme]
	return
}

func AvailableThemes() []string {
	list := make([]string, 0, len(themes))

	for t := range themes {
		list = append(list, t)
	}
	return list
}
