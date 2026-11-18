package memory

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
)

type Repository struct {
	themedBadges map[string]string
	badgeNames   map[string]struct{}
}

func New(dir string) (*Repository, error) {
	badges := make(map[string]string)
	names := make(map[string]struct{})

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".svg" {
			return err
		}

		svg, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}
		svgStr := string(svg)

		tagStart := strings.Index(svgStr, "<svg")
		if tagStart == -1 {
			return nil
		}

		tagOffset := strings.Index(svgStr[tagStart:], ">")
		if tagOffset == -1 {
			return nil
		}

		start := tagStart + tagOffset + 1
		end := strings.LastIndex(svgStr, "</svg>")
		if end == -1 || start >= end {
			return nil
		}

		content := svgStr[start:end]

		themedName := strings.ToLower(strings.TrimSuffix(filepath.Base(path), ".svg"))
		badges[themedName] = strings.TrimSpace(content)

		name := domain.TrimThemeSuffix(themedName)
		names[name] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	return &Repository{themedBadges: badges, badgeNames: names}, nil
}

func (r *Repository) Get(_ context.Context, name, theme string) (string, error) {
	if _, ok := r.badgeNames[name]; !ok {
		return "", domain.ErrBadgeNotFound
	}

	themeKey := name + "-" + theme
	if badge, ok := r.themedBadges[themeKey]; ok {
		return badge, nil
	}

	if badge, ok := r.themedBadges[name]; ok {
		return badge, nil
	}

	return "", domain.ErrBadgeNotFound
}
