package diff

import "strings"

type stringInterface struct {
	a []string
	b []string
}

func (si stringInterface) LenA() int {
	return len(si.a)
}

func (si stringInterface) LenB() int {
	return len(si.b)
}

func (si stringInterface) Compare(i, j int) int {
	return strings.Compare(si.a[i], si.b[j])
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
