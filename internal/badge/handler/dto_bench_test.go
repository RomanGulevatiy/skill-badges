package handler

import (
	"strings"
	"testing"
)

func BenchmarkGetBadgesRequestToInput(b *testing.B) {
	maxCountBadges := make([]string, 0, maxBadges)
	for range maxBadges {
		maxCountBadges = append(maxCountBadges, "go")
	}

	tests := []struct {
		name string
		req  GetBadgesRequest
	}{
		{
			name: "typical",
			req: GetBadgesRequest{
				Badges:  "js,ts,vite,vue,pinia,tailwind",
				Perline: "15",
				Theme:   "dark",
				Size:    "md",
			},
		},
		{
			name: "aliases and whitespace",
			req: GetBadgesRequest{
				Badges:  " JS, tS , go,py  , Rs ",
				Perline: "3",
				Theme:   "light",
				Size:    "xl",
			},
		},
		{
			name: "max badges",
			req: GetBadgesRequest{
				Badges:  strings.Join(maxCountBadges, ","),
				Perline: "15",
				Theme:   "dark",
				Size:    "md",
			},
		},
		{
			name: "longest allowed param",
			req: GetBadgesRequest{
				Badges:  strings.Repeat("c", maxBadgesParamLength),
				Perline: "15",
				Theme:   "dark",
				Size:    "md",
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := test.req.ToInput(); err != nil {
					b.Fatalf("ToInput() error: %v", err)
				}
			}
		})
	}
}

func TestToInputAllocations(t *testing.T) {
	tests := []struct {
		name string
		req  GetBadgesRequest
	}{
		{
			name: "typical",
			req: GetBadgesRequest{
				Badges:  "js,ts,vite,vue,pinia,tailwind",
				Perline: "15",
				Theme:   "dark",
				Size:    "md",
			},
		},
		{
			name: "max badges",
			req: GetBadgesRequest{
				Badges:  strings.Repeat("go,", maxBadges-1) + "go",
				Perline: "15",
				Theme:   "dark",
				Size:    "md",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			allocs := testing.AllocsPerRun(100, func() {
				in, err := test.req.ToInput()
				if err != nil {
					t.Fatalf("ToInput() error: %v", err)
				}
				_ = in
			})
			if allocs > 2 {
				t.Fatalf("ToInput() allocations = %.0f, want <= 2", allocs)
			}
		})
	}
}
