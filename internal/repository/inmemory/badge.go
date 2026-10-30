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
	badges map[string]string
	logger *slog.Logger
}

func NewBadgeInmemoryRepo(dir string, log *slog.Logger) (*BadgeInmemoryRepo, error) {
	badges := make(map[string]string)

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

		name := strings.TrimSuffix(filepath.Base(path), ".svg")
		badges[strings.ToLower(name)] = strings.TrimSpace(content)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	return &BadgeInmemoryRepo{badges: badges, logger: log}, nil
}

func (r *BadgeInmemoryRepo) GetBadge(ctx context.Context, name string) (string, error) {
	lowerName := strings.ToLower(name)

	badge, ok := r.badges[lowerName]
	if !ok {
		r.logger.DebugContext(ctx, "badge not found in repository", slog.String("name", lowerName))
		return "", domain.ErrBadgeNotFound
	}
	return badge, nil
}
