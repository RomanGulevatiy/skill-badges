package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
	"github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
)

type mockRepository struct {
	mockGet func(ctx context.Context, name, theme string) (string, error)
}

func (r *mockRepository) Get(ctx context.Context, name, theme string) (string, error) {
	if r.mockGet != nil {
		return r.mockGet(ctx, name, theme)
	}
	return "", domain.ErrBadgeNotFound
}

func TestUsecase_GetBadges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          usecase.GetBadgesInput
		mockGet     func(ctx context.Context, name, theme string) (string, error)
		expectedSVG string
		expectedErr error
	}{
		{
			name: "returns badges in a single row",
			in: usecase.GetBadgesInput{
				Badges:  []string{"golang", "mongodb"},
				Perline: 3,
				Theme:   "dark",
				Dim:     domain.Dimensions{Size: 48, Gap: 8},
			},
			mockGet: func(ctx context.Context, name, theme string) (string, error) {
				badges := map[string]string{
					"golang":  `<path id="golang"/>`,
					"mongodb": `<path id="mongodb"/>`,
				}

				if theme != "dark" {
					t.Errorf("expected theme to be dark, got %s", theme)
				}

				b, ok := badges[name]
				if !ok {
					return "", domain.ErrBadgeNotFound
				}
				return b, nil
			},
			expectedSVG: `<svg xmlns="http://www.w3.org/2000/svg" width="104" height="48" fill="none" viewBox="0 0 104 48">` +
				`<svg x="0" y="0" width="48" height="48" viewBox="0 0 256 256"><path id="golang"/></svg>` +
				`<svg x="56" y="0" width="48" height="48" viewBox="0 0 256 256"><path id="mongodb"/></svg>` +
				`</svg>`,
			expectedErr: nil,
		},
		{
			name: "wraps badges across multiple rows",
			in: usecase.GetBadgesInput{
				Badges:  []string{"golang", "mongodb", "docker"},
				Perline: 2,
				Theme:   "light",
				Dim:     domain.Dimensions{Size: 48, Gap: 8},
			},
			mockGet: func(ctx context.Context, name, theme string) (string, error) {
				badges := map[string]string{
					"golang":  `<path id="golang"/>`,
					"mongodb": `<path id="mongodb"/>`,
					"docker":  `<path id="docker"/>`,
				}

				if theme != "light" {
					t.Errorf("expected theme to be light, got %s", theme)
				}

				b, ok := badges[name]
				if !ok {
					return "", domain.ErrBadgeNotFound
				}
				return b, nil
			},
			expectedSVG: `<svg xmlns="http://www.w3.org/2000/svg" width="104" height="104" fill="none" viewBox="0 0 104 104">` +
				`<svg x="0" y="0" width="48" height="48" viewBox="0 0 256 256"><path id="golang"/></svg>` +
				`<svg x="56" y="0" width="48" height="48" viewBox="0 0 256 256"><path id="mongodb"/></svg>` +
				`<svg x="0" y="56" width="48" height="48" viewBox="0 0 256 256"><path id="docker"/></svg>` +
				`</svg>`,
			expectedErr: nil,
		},
		{
			name: "skips badges that are not found",
			in: usecase.GetBadgesInput{
				Badges:  []string{"golang", "unknown", "docker"},
				Perline: 15,
				Theme:   "light",
				Dim:     domain.Dimensions{Size: 48, Gap: 8},
			},
			mockGet: func(ctx context.Context, name, theme string) (string, error) {
				badges := map[string]string{
					"golang": `<path id="golang"/>`,
					"docker": `<path id="docker"/>`,
				}

				if theme != "light" {
					t.Errorf("expected theme to be light, got %s", theme)
				}

				b, ok := badges[name]
				if !ok {
					return "", domain.ErrBadgeNotFound
				}
				return b, nil
			},
			expectedSVG: `<svg xmlns="http://www.w3.org/2000/svg" width="104" height="48" fill="none" viewBox="0 0 104 48">` +
				`<svg x="0" y="0" width="48" height="48" viewBox="0 0 256 256"><path id="golang"/></svg>` +
				`<svg x="56" y="0" width="48" height="48" viewBox="0 0 256 256"><path id="docker"/></svg>` +
				`</svg>`,
			expectedErr: nil,
		},
		{
			name: "returns an error when no badges are found",
			in: usecase.GetBadgesInput{
				Badges:  []string{"unknown1", "unknown2"},
				Perline: 15,
				Theme:   "dark",
				Dim:     domain.Dimensions{Size: 48, Gap: 8},
			},
			mockGet: func(ctx context.Context, name, theme string) (string, error) {
				return "", domain.ErrBadgeNotFound
			},
			expectedSVG: "",
			expectedErr: domain.ErrBadgeNotFound,
		},
		{
			name: "returns an error without querying the repository for empty input",
			in: usecase.GetBadgesInput{
				Badges:  []string{},
				Perline: 15,
				Theme:   "dark",
				Dim:     domain.Dimensions{Size: 48, Gap: 8},
			},
			mockGet: func(ctx context.Context, name, theme string) (string, error) {
				t.Fatal("badge.repository.memory.Get must not be called for empty input")
				return "", nil
			},
			expectedSVG: "",
			expectedErr: domain.ErrBadgeNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := &mockRepository{mockGet: test.mockGet}
			uc := usecase.New(repo)

			gotSVG, err := uc.GetBadges(t.Context(), test.in)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("expected error %v, got %v", test.expectedErr, err)
			}

			if gotSVG != test.expectedSVG {
				t.Errorf("SVG mismatch:\ngot:  %s\nwant: %s", gotSVG, test.expectedSVG)
			}
		})
	}
}
