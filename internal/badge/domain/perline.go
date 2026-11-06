package domain

const (
	defaultPerline = 15
	minPerline     = 1
	maxPerline     = 32
)

func ValidatePerline(perline int) int {
	if perline < minPerline || perline > maxPerline {
		perline = defaultPerline
	}
	return perline
}
