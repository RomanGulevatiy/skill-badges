package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/domain"
)

type BadgeInmemoryRepo struct {
	themedBadges map[string]string
	badgeNames   map[string]struct{}
	logger       *slog.Logger
}

func NewBadgeInmemoryRepo(dir string, log *slog.Logger) (*BadgeInmemoryRepo, error) {
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

		content := strings.TrimSpace(svgStr[start:end])

		themedName := strings.TrimSuffix(filepath.Base(path), ".svg")
		lowerThemedName := strings.ToLower(themedName)

		badges[lowerThemedName] = strings.TrimSpace(content)

		name := lowerThemedName
		for theme := range domain.Themes {
			if before, found := strings.CutSuffix(name, "-"+theme); found {
				name = before
				break
			}
		}

		names[name] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	return &BadgeInmemoryRepo{themedBadges: badges, badgeNames: names, logger: log}, nil
}

func (r *BadgeInmemoryRepo) GetBadge(ctx context.Context, name, theme string) (string, error) {
	lowerName := strings.ToLower(name)

	if _, ok := r.badgeNames[lowerName]; !ok {
		r.logger.DebugContext(ctx, "invalid or explicitly suffixed badge name requested", slog.String("name", lowerName))
		return "", domain.ErrBadgeNotFound
	}

	themeKey := lowerName + "-" + theme
	if badge, ok := r.themedBadges[themeKey]; ok {
		return badge, nil
	}

	if badge, ok := r.themedBadges[lowerName]; ok {
		return badge, nil
	}

	r.logger.DebugContext(ctx, "badge not found in repository", slog.String("name", lowerName))
	return "", domain.ErrBadgeNotFound
}
