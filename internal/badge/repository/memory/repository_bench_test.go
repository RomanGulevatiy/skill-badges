package memory

import (
	"context"
	"testing"
)

func BenchmarkRepositoryGet(b *testing.B) {
	const content = `<path fill="#fff" d="m136.65"/>`

	names := []string{
		"golang", "javascript", "typescript", "docker", "git", "github",
		"linux", "debian", "mongodb", "postgresql", "redis", "vite",
	}

	noThemeVariant := map[string]struct{}{
		"docker": {},
		"git":    {},
	}

	themedBadges := make(map[string]string, len(names)*2)
	badgeNames := make(map[string]struct{}, len(names))
	for _, name := range names {
		badgeNames[name] = struct{}{}
		themedBadges[name] = content
		if _, ok := noThemeVariant[name]; !ok {
			themedBadges[name+"-dark"] = content
		}
	}

	repo := &Repository{themedBadges: themedBadges, badgeNames: badgeNames}
	ctx := context.Background()

	tests := []struct {
		name  string
		badge string
	}{
		{name: "themed variant", badge: "javascript"},
		{name: "fallback without theme", badge: "docker"},
		{name: "missing", badge: "unknown"},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = repo.Get(ctx, test.badge, "dark")
			}
		})
	}
}
