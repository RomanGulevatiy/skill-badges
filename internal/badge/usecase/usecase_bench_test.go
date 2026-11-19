package usecase

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
)

type benchRepository map[string]string

func (r benchRepository) Get(_ context.Context, name, _ string) (string, error) {
	badge, ok := r[name]
	if !ok {
		return "", domain.ErrBadgeNotFound
	}
	return badge, nil
}

func benchBadgeContent(size int) string {
	const prefix = `<path fill="#fff" d="`
	const suffix = `"/>`

	repeats := max((size-len(prefix)-len(suffix))/len("m136.65"), 1)

	var b strings.Builder
	b.Grow(len(prefix) + repeats*len("m136.65") + len(suffix))
	b.WriteString(prefix)
	for range repeats {
		b.WriteString("m136.65")
	}
	b.WriteString(suffix)
	return b.String()
}

func benchBadges(names []string, content string) benchRepository {
	repo := make(benchRepository, len(names))
	for _, name := range names {
		repo[name] = content
	}
	return repo
}

func benchmarkBadgeNames(count int) []string {
	names := make([]string, count)
	for i := range names {
		names[i] = "badge" + strconv.Itoa(i)
	}
	return names
}

func BenchmarkUsecaseGetBadges(b *testing.B) {
	const averageAssetSize = 1800
	content := benchBadgeContent(averageAssetSize)

	tests := []struct {
		name    string
		count   int
		perline int
	}{
		{name: "typical 8 badges", count: 8, perline: 15},
		{name: "full row 15 badges", count: 15, perline: 15},
		{name: "max 100 badges", count: 100, perline: 15},
		{name: "max 100 badges with perline 32", count: 100, perline: 32},
	}

	for _, test := range tests {
		repo := benchBadges(benchmarkBadgeNames(test.count), content)
		in := GetBadgesInput{
			Badges:  benchmarkBadgeNames(test.count),
			Perline: test.perline,
			Theme:   "dark",
			Dim:     domain.Dimensions{Size: 48, Gap: 8},
		}
		uc := New(repo)

		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := uc.GetBadges(context.Background(), in); err != nil {
					b.Fatalf("GetBadges() error: %v", err)
				}
			}
		})
	}
}

func BenchmarkRender(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{name: "small 320B", size: 320},
		{name: "average 1800B", size: 1800},
		{name: "large 8000B", size: 8000},
	}

	counts := []struct {
		name    string
		count   int
		perline int
	}{
		{name: "8 badges", count: 8, perline: 15},
		{name: "100 badges", count: 100, perline: 15},
		{name: "100 badges with perline 32", count: 100, perline: 32},
	}

	dim := domain.Dimensions{Size: 48, Gap: 8}

	for _, size := range sizes {
		content := benchBadgeContent(size.size)
		for _, count := range counts {
			badges := make([]string, count.count)
			for i := range badges {
				badges[i] = content
			}

			b.Run(size.name+"/"+count.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					render(badges, count.perline, dim)
				}
			})
		}
	}
}
