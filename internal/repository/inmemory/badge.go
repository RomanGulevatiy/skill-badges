package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/domain"
)

var contentRegex = regexp.MustCompile(`(?s)<svg[^>]*>(.*?)</svg>`)

type BadgeInmemoryRepo struct {
	badges map[string]string
	logger *slog.Logger
}

func NewBadgeInmemoryRepo(dir string, logger *slog.Logger) (*BadgeInmemoryRepo, error) {
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

		var content string
		if matches := contentRegex.FindStringSubmatch(svgStr); len(matches) > 1 {
			content = strings.TrimSpace(matches[1])
		} else {
			return nil
		}

		name := strings.TrimSuffix(filepath.Base(path), ".svg")
		badges[strings.ToLower(name)] = strings.TrimSpace(content)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	return &BadgeInmemoryRepo{badges: badges, logger: logger}, nil
}

func (r *BadgeInmemoryRepo) GetBadge(ctx context.Context, name string) (string, error) {
	lowerName := strings.ToLower(name)

	badge, ok := r.badges[lowerName]
	if !ok {
		r.logger.WarnContext(ctx, "badge not found in repository", slog.String("name", lowerName))
		return "", domain.ErrBadgeNotFound
	}
	return badge, nil
}
