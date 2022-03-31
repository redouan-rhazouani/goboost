package chans

import (
	"testing"
)

func TestMerge(t *testing.T) {
	n := 20
	c1 := Repeat(1, n)
	c2 := Repeat(2, n)
	a := make([]int, 0, n)
	for v := range Merge(c1, c2) {
		a = append(a, v)
	}
	if m := len(a); m != 2*n {
		t.Errorf("Number of elements got=%d want=%d", m, 2*n)
	}
}
