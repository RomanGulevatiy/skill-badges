package handler_test

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
	"github.com/RomanGulevatiy/skill-badges/internal/badge/handler"
	"github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
)

type mockUsecase struct {
	mockGetBadges func(ctx context.Context, in usecase.GetBadgesInput) (string, error)
}

func (u *mockUsecase) GetBadges(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
	if u.mockGetBadges != nil {
		return u.mockGetBadges(ctx, in)
	}
	return "", nil
}

func TestHandler_GetBadges(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.DiscardHandler)

	tests := []struct {
		name                string
		url                 string
		mockGetBadges       func(ctx context.Context, in usecase.GetBadgesInput) (string, error)
		expectedStatus      int
		expectedBody        string
		expectedContentType string
	}{
		{
			name: "returns svg with parsed query params",
			url:  "/api/v1/badges?b=go,git,docker&perline=3&theme=light&size=xl",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				expectedInput := usecase.GetBadgesInput{
					Badges:  []string{"golang", "git", "docker"},
					Perline: 3,
					Theme:   "light",
					Dim:     domain.Dimensions{Size: 80, Gap: 16},
				}

				if !reflect.DeepEqual(in, expectedInput) {
					t.Errorf("unexpected input passed to usecase: got %+v, want %+v", in, expectedInput)
				}
				return "<svg></svg>", nil
			},
			expectedStatus:      http.StatusOK,
			expectedBody:        `<svg></svg>`,
			expectedContentType: "image/svg+xml",
		},
		{
			name: "no badge parameter provided",
			url:  "/api/v1/badges",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				t.Fatal("badge.usecase.GetBadges must not be called when validation fails")
				return "", nil
			},
			expectedStatus:      http.StatusBadRequest,
			expectedBody:        "at least one badge must be specified in 'b'\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name: "falls back to default perline when perline is zero",
			url:  "/api/v1/badges?b=go,docker&perline=0",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				expectedInput := usecase.GetBadgesInput{
					Badges:  []string{"golang", "docker"},
					Perline: 15,
					Theme:   "dark",
					Dim:     domain.Dimensions{Size: 48, Gap: 8},
				}

				if !reflect.DeepEqual(in, expectedInput) {
					t.Errorf("unexpected input passed to usecase: got %+v, want %+v", in, expectedInput)
				}
				return "<svg></svg>", nil
			},
			expectedStatus:      http.StatusOK,
			expectedBody:        `<svg></svg>`,
			expectedContentType: "image/svg+xml",
		},
		{
			name: "falls back to default perline when perline is non-numeric",
			url:  "/api/v1/badges?b=go,docker,mongodb&perline=ABC",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				expectedInput := usecase.GetBadgesInput{
					Badges:  []string{"golang", "docker", "mongodb"},
					Perline: 15,
					Theme:   "dark",
					Dim:     domain.Dimensions{Size: 48, Gap: 8},
				}

				if !reflect.DeepEqual(in, expectedInput) {
					t.Errorf("unexpected input passed to usecase: got %+v, want %+v", in, expectedInput)
				}
				return "<svg></svg>", nil
			},
			expectedStatus:      http.StatusOK,
			expectedBody:        `<svg></svg>`,
			expectedContentType: "image/svg+xml",
		},
		{
			name: "falls back to default theme when theme is unsupported",
			url:  "/api/v1/badges?b=git,docker&theme=catppuccin",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				expectedInput := usecase.GetBadgesInput{
					Badges:  []string{"git", "docker"},
					Perline: 15,
					Theme:   "dark",
					Dim:     domain.Dimensions{Size: 48, Gap: 8},
				}

				if !reflect.DeepEqual(in, expectedInput) {
					t.Errorf("unexpected input passed to usecase: got %+v, want %+v", in, expectedInput)
				}
				return "<svg></svg>", nil
			},
			expectedStatus:      http.StatusOK,
			expectedBody:        `<svg></svg>`,
			expectedContentType: "image/svg+xml",
		},
		{
			name: "falls back to default size when size is unsupported",
			url:  "/api/v1/badges?b=linux,arch,debian&size=3xl",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				expectedInput := usecase.GetBadgesInput{
					Badges:  []string{"linux", "archlinux", "debian"},
					Perline: 15,
					Theme:   "dark",
					Dim:     domain.Dimensions{Size: 48, Gap: 8},
				}

				if !reflect.DeepEqual(in, expectedInput) {
					t.Errorf("unexpected input passed to usecase: got %+v, want %+v", in, expectedInput)
				}
				return "<svg></svg>", nil
			},
			expectedStatus:      http.StatusOK,
			expectedBody:        `<svg></svg>`,
			expectedContentType: "image/svg+xml",
		},
		{
			name: "badges not found",
			url:  "/api/v1/badges?b=unknown",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				return "", domain.ErrBadgeNotFound
			},
			expectedStatus:      http.StatusNotFound,
			expectedBody:        domain.ErrBadgeNotFound.Error() + "\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name: "internal server error",
			url:  "/api/v1/badges?b=go",
			mockGetBadges: func(ctx context.Context, in usecase.GetBadgesInput) (string, error) {
				return "", errors.New("some error")
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        "internal server error\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			uc := &mockUsecase{mockGetBadges: test.mockGetBadges}
			h := handler.New(uc, logger)

			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			rr := httptest.NewRecorder()

			h.GetBadges(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, rr.Code)
			}

			if gotBody := rr.Body.String(); gotBody != test.expectedBody {
				t.Errorf("expected body %q, got %q", test.expectedBody, gotBody)
			}

			if ct := rr.Header().Get("Content-Type"); ct != test.expectedContentType {
				t.Errorf("expected Content-Type %q, got %q", test.expectedContentType, ct)
			}

			if test.expectedStatus == http.StatusOK {
				if cc := rr.Header().Get("Cache-Control"); cc != "public, max-age=86400" {
					t.Errorf("unexpected Cache-Control header: %s", cc)
				}
			}
		})
	}
}
