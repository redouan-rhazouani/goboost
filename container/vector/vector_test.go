package vector

import (
	"testing"
)

func eq(v int) func(int) bool {
	return func(x int) bool { return x == v }
}

func TestVector(t *testing.T) {
	vec := Make[int](0, 3)
	if m, n := vec.Len(), vec.Cap(); m != 0 || n != 3 {
		t.Errorf("(len, cap) got = (%d,%d )want (%d,%d)", m, n, 0, 3)
	}
	vec = Make[int](1, 3)
	if m, n := vec.Len(), vec.Cap(); m != 1 || n != 3 {
		t.Errorf("(len, cap) got = (%d,%d )want (%d,%d)", m, n, 1, 3)
	}
}

func TestSwapDelete(t *testing.T) {
	tests := []struct {
		xs   []int
		idx  int
		ys   []int
		elem int
	}{
		{[]int{1}, 0, []int{}, 1},
		{[]int{1, 2}, 0, []int{2}, 1},
		{[]int{1, 2}, 1, []int{1}, 2},

		{[]int{1, 2, 3}, 0, []int{3, 2}, 1},
		{[]int{1, 2, 3}, 1, []int{1, 3}, 2},
		{[]int{1, 2, 3}, 2, []int{1, 2}, 3},
	}

	for i, tc := range tests {
		v := Vector[int](tc.xs)
		v = v.Copy()
		s := v
		if e := v.SwapDelete(tc.idx); e != tc.elem {
			t.Errorf("case-%d: %v.SwapDelete(%d) got=%v want %v", i, e, tc.idx, e, tc.elem)

		}
		if !Equal(v, Vector[int](tc.ys)) {
			t.Errorf("case-%d: %v.SwapDelete(%d) got=%v want %v", i, tc.xs, tc.idx, v, tc.ys)
		}

		if n := len(s); n > 0 && s[n-1] != 0 {
			t.Fatalf("last element in original slice must be zero got = %v", s[n-1])
		}

	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		xs   []int
		idx  int
		ys   []int
		elem int
	}{
		{[]int{1}, 0, []int{}, 1},
		{[]int{1, 2}, 0, []int{2}, 1},
		{[]int{1, 2}, 1, []int{1}, 2},

		{[]int{1, 2, 3, 4, 5}, 0, []int{2, 3, 4, 5}, 1},
		{[]int{1, 2, 3, 4, 5}, 1, []int{1, 3, 4, 5}, 2},
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}, 3},
		{[]int{1, 2, 3, 4, 5}, 3, []int{1, 2, 3, 5}, 4},
		{[]int{1, 2, 3, 4, 5}, 4, []int{1, 2, 3, 4}, 5},
	}

	for i, tc := range tests {
		v := Vector[int](tc.xs)
		v = v.Copy()
		if e := v.Delete(tc.idx); e != tc.elem {
			t.Errorf("case-%d: %v.Delete(%d) got=%v want %v", i, e, tc.idx, e, tc.elem)

		}
		if !Equal(v, Vector[int](tc.ys)) {
			t.Errorf("case-%d: %v.Delete(%d) got=%v want %v", i, tc.xs, tc.idx, v, tc.ys)
		}
	}
}

func TestDeleteRange(t *testing.T) {
	tests := []struct {
		xs        []int
		low, high int
		ys        []int
	}{
		{nil, 0, 0, nil},
		{[]int{1}, 0, 1, []int{}},
		{[]int{1, 2}, 0, 2, []int{}},
		{[]int{1, 2}, 1, 2, []int{1}},

		{[]int{1, 2, 3, 4, 5}, 0, 2, []int{3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 1, 4, []int{1, 5}},
		{[]int{1, 2, 3, 4, 5}, 2, 5, []int{1, 2}},
		{[]int{1, 2, 3, 4, 5}, 0, 5, []int{}},
	}

	for _, c := range tests {
		v := Vector[int](c.xs).Copy()
		os := v
		v.DeleteRange(c.low, c.high)
		if !Equal(v, c.ys) {
			t.Errorf("%v.DeleteRange(%d, %d) got=%v want %v", c.xs, c.low, c.high, v, c.ys)

		}
		if m, n := Count(os, 0), len(c.xs)-v.Len(); m != n {
			t.Errorf("Count(%v, 0) got=%v want %v", os, m, n)
		}
		v = Vector[int](c.xs).Copy()
		os = v
		v.Clear()
		if n := v.Len(); n != 0 {
			t.Errorf("%v.Clear().Len() got %v want 0", c.xs, n)
		}
		if m, n := Count(os, 0), len(c.xs); m != n {
			t.Errorf("Count(%v, 0) got=%v want %v", os, m, n)
		}
	}
}

func TestAt(t *testing.T) {
	v := Vector[int]([]int{1, 2, 3, 4, 5})

	for i, x := range v {
		if v.At(i) != x {
			t.Errorf("%v.At(%v) got = %d want = %d", v, i, v.At(i), x)
		}
	}
	u := Vector[int]([]int{5, 4, 3, 2, 1})
	for i, j := 0, -1; j >= -len(v); i, j = i+1, j-1 {
		want := u[i]
		if v.At(j) != u[i] {
			t.Errorf("%v.At(%v) got = %d want = %d", v, j, v.At(j), want)
		}
	}
}

func TestIndex(t *testing.T) {

	tests := []struct {
		vec []int
		v   int
		i   int
	}{
		{[]int{}, 1, -1},
		{[]int{1}, 1, 0},
		{[]int{0}, 1, -1},
		{[]int{1, 2, 3}, 0, -1},
		{[]int{1, 2, 2, 3}, 2, 1},
		{[]int{1, 2, 3, 3}, 3, 2},
		{[]int{1, 2, 3}, 4, -1},
	}

	for _, c := range tests {
		vec := Vector[int](c.vec)
		if g := Index(vec, c.v); g != c.i {
			t.Errorf("%v.Index(%v) got = %d want = %d", vec, c.v, g, c.i)
		}
		if g := vec.IndexFunc(eq(c.v)); g != c.i {
			t.Errorf("%v.IndexFunc(%v) got = %d want = %d", vec, c.v, g, c.i)
		}

	}
}

func TestLastIndex(t *testing.T) {

	tests := []struct {
		vec []int
		v   int
		i   int
	}{
		{[]int{}, 1, -1},
		{[]int{1}, 1, 0},
		{[]int{0}, 1, -1},
		{[]int{1, 2, 3}, 0, -1},
		{[]int{1, 2, 2, 3}, 2, 2},
		{[]int{1, 2, 3, 3}, 3, 3},
		{[]int{1, 2, 3}, 4, -1},
	}

	eq := func(v int) func(int) bool {
		return func(x int) bool { return x == v }
	}

	for _, c := range tests {
		vec := Vector[int](c.vec)
		if g := LastIndex(vec, c.v); g != c.i {
			t.Errorf("%v.LastIndex(%v) got = %d want = %d", vec, c.v, g, c.i)
		}
		if g := vec.LastIndexFunc(eq(c.v)); g != c.i {
			t.Errorf("%v.LastIndexFunc(%v) got = %d want = %d", vec, c.v, g, c.i)
		}

	}
}

func TestEqual(t *testing.T) {
	v1 := Make[int](0, 7)
	if !Equal(v1, v1) {
		t.Errorf("Empty vector must be equal to itself")
	}
	v1.Append(1, 2, 3)
	if !Equal(v1, v1) {
		t.Errorf("Empty vector must be equal to itself")
	}
	u := v1.Copy()
	if !Equal(v1, u) {
		t.Errorf("Empty vector must be equal to its copy")
	}

	if v2 := v1[2:]; Equal(v1, v2) {
		t.Errorf("%v.Equal(%v) got = %v want %v", v1, v2, true, false)
	}
	if v2 := v1[1:]; Equal(v1, v2) {
		t.Errorf("%v.Equal(%v) got = %v want %v", v1, v2, true, false)
	}
	v2 := v1.Copy()
	v2[len(v2)-1] = 11
	if Equal(v1, v2) {
		t.Errorf("%v.Equal(%v) got = %v want %v", v1, v2, true, false)
	}
}

func TestAppend(t *testing.T) {
	v1 := Make[int](0, 7)
	v1.Append(1, 2)
	v1.Push(3)
	v1.Push(4)
	v1.Append(5)
	v1.Pop()
	v1.Append(5)
	v1.Push(6)
	v2 := []int{1, 2, 3, 4, 5, 6}
	if !Equal(v1, v2) {
		t.Errorf("append got %v want %v", v1, v2)
	}
	v3 := Make[int](0, 7)
	for v1.Len() > 0 {
		v3.Push(v1.Pop())
	}
	v2 = []int{1, 2, 3, 4, 5, 6}
	if !Equal(v3, []int{6, 5, 4, 3, 2, 1}) {
		t.Errorf("append got %v want %v", v1, v2)
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		xs  []int
		val int
		ys  []int
	}{
		{nil, 0, nil},
		{[]int{2}, 1, []int{2}},
		{[]int{1, 2}, 3, []int{1, 2}},
		{[]int{2, 1}, 1, []int{2}},
		{[]int{1, 2}, 1, []int{2}},
		{[]int{1, 1}, 1, []int{}},

		{[]int{1, 2, 1, 3, 1}, 1, []int{2, 3}},
		{[]int{2, 1, 1, 4, 1}, 1, []int{2, 4}},
	}

	for _, c := range tests {
		v := Vector[int](c.xs).Copy()
		os := v
		v.RemoveFunc(eq(c.val))
		if !Equal(v, c.ys) {
			t.Errorf("%v.RemoveFunc(%d) got=%v want %v", c.xs, c.val, v, c.ys)

		}
		if m, n := Count(os, 0), len(c.xs)-v.Len(); m != n {
			t.Errorf("Count(%v, 0) got=%v want %v", os, m, n)
		}
	}

	for _, c := range tests {
		v := Vector[int](c.xs).Copy()
		os := v
		Remove(&v, c.val)
		if !Equal(v, c.ys) {
			t.Errorf("Remove(%v,%d) got=%v want %v", c.xs, c.val, v, c.ys)

		}
		if m, n := Count(os, 0), len(c.xs)-v.Len(); m != n {
			t.Errorf("Count(%v, 0) got=%v want %v", os, m, n)
		}
	}
}
