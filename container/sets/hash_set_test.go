package sets

import (
	"testing"
)

func checkSetLen[T comparable](t *testing.T, s USet[T], len int) bool {
	t.Helper()
	if n := s.Len(); n != len {
		t.Errorf("s.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkSet[T comparable](t *testing.T, s USet[T], es []T) bool {
	t.Helper()
	if !checkSetLen(t, s, len(es)) {
		return false
	}

	for _, v := range es {
		if !s.Contains(v) {
			t.Errorf("!s.Contains(%v)", v)
			return false
		}
	}
	return true
}

func TestSetLen(t *testing.T) {
	s := Make[int]()
	checkSetLen(t, s, 0)
	s.Add(1)
	s.Add(2)
	s.Add(1)
	s.InsertSlice(nil)
	s.InsertSlice([]int{})
	s.InsertSlice([]int{1, 2, 3})
	s.InsertSlice([]int{3, 4, 5})
	l := []int{1, 2, 3, 4, 5}
	checkSet(t, s, l)
	s.Delete(2)
	s.Delete(4)
	checkSet(t, s, []int{1, 3, 5})
	s.DeleteIF(func(v int) bool { return v&1 == 1 })
	checkSetLen(t, s, 0)
	s.Update(FromSlice[int](nil))
	s.Update(FromSlice([]int{}))
	s.Update(FromSlice([]int{1}))
	s.Update(FromSlice([]int{1, 2, 3}))
	s.Update(FromSlice([]int{4, 4}))
	s.Update(FromSlice([]int{1, 4, 5}))
	checkSet(t, s, l)
	s.Clear()
	checkSetLen(t, s, 0)
}

func TestDisjoint(t *testing.T) {
	emptySet := Make[int]()
	tests := []struct {
		xs, ys USet[int]
		want   bool
	}{
		{emptySet, emptySet, true},
		{nil, nil, true},
		{USet[int]{}, USet[int]{}, true},
		{FromSlice([]int{1}), FromSlice([]int{2}), true},
		{FromSlice([]int{1, 3}), FromSlice([]int{2, 4}), true},
		{FromSlice([]int{1, 2, 3}), FromSlice([]int{-1, 2, 5}), false},
	}
	for _, tc := range tests {
		s1, s2 := tc.xs, tc.ys
		if g := s1.IsDisjoint(s2); g != tc.want {
			t.Errorf("%v.disjoint(%v) got = %v want = %v", s1, s2, g, tc.want)
		}
		if g := s2.IsDisjoint(s1); g != tc.want {
			t.Errorf("%v.disjoint(%v) got = %v want = %v", s1, s2, g, tc.want)
		}
	}
}

func TestSubsetSuperSet(t *testing.T) {
	xs := FromSlice([]int{1, 2, 3})
	tests := []struct {
		xs, ys     USet[int]
		issubset   bool
		issuperset bool
	}{
		{xs, xs, true, true},
		{nil, nil, true, true},
		{USet[int]{}, USet[int]{}, true, true},
		{FromSlice([]int{1}), FromSlice([]int{1}), true, true},
		{FromSlice([]int{}), FromSlice([]int{1}), true, false},
		{FromSlice[int](nil), FromSlice([]int{1}), true, false},
		{FromSlice([]int{1}), FromSlice([]int{2}), false, false},
		{FromSlice([]int{1}), FromSlice([]int{1, 2}), true, false},
		{FromSlice([]int{0, 1, 2}), FromSlice([]int{0, 100, 231, 2}), false, false},
		{FromSlice([]int{0, 1, 2}), FromSlice([]int{0, 100, 1, 231, 2}), true, false},
	}
	for _, tc := range tests {
		s1, s2 := tc.xs, tc.ys
		if g := s1.IsSubset(s2); g != tc.issubset {
			t.Errorf("%v.issubset(%v) got = %v want = %v", s1, s2, g, tc.issubset)
		}
		if tc.issubset {
			if g := s2.IsSuperset(s1); g != tc.issubset {
				t.Errorf("%v.IsSuperset(%v) got = %v want = %v", s2, s1, g, tc.issubset)
			}
		}
		if g := s1.IsSuperset(s2); g != tc.issuperset {
			t.Errorf("%v.IsSuperset(%v) got = %v want = %v", s1, s2, g, tc.issubset)
		}
		if tc.issuperset {
			if g := s2.IsSubset(s1); g != tc.issubset {
				t.Errorf("%v.IsSubset(%v) got = %v want = %v", s2, s1, g, tc.issubset)
			}
		}

	}
}

func TestIntersection(t *testing.T) {
	xs := []int{3, 5, 11, 77}
	tests := []struct {
		xs, ys []int
		want   []int
	}{
		{nil, nil, nil},
		{xs, xs, xs},
		{[]int{}, []int{}, []int{}},
		{xs, []int{11}, []int{11}},
		{xs, []int{11, 3}, []int{11, 3}},
		{xs, []int{}, []int{}},
		{xs, nil, []int{}},
		{[]int{11, 1, 3, 77, 103, 5, -5}, []int{2, 11, 77, -7, -23, 5, 3}, xs},
		{[]int{11, 1, 3, 77, 103, 5, -5, 9}, []int{2, 11, 77, -7, -23, 5, 3}, xs},
	}
	for _, tc := range tests {
		s1, s2 := FromSlice(tc.xs), FromSlice(tc.ys)
		s1.IntersectionUpdate(s2)
		checkSet(t, s1, tc.want)
	}
	for _, tc := range tests {
		s1, s2 := FromSlice(tc.xs), FromSlice(tc.ys)
		checkSet(t, Intersection(s1, s2), tc.want)
	}
}

func TestDifference(t *testing.T) {
	xs := []int{3, 5, 11, 77}
	ys := []int{1, 2, 6, 12, 43}
	zs := append(ys, xs...)
	tests := []struct {
		xs, ys []int
		want   []int
	}{
		{nil, nil, nil},
		{xs, xs, nil},
		{[]int{}, []int{}, []int{}},
		{xs, []int{3}, xs[1:]},
		{xs, xs[:2], xs[2:]},
		{zs, xs, ys},
		{zs, ys, xs},
	}
	for _, tc := range tests {
		s1, s2 := FromSlice(tc.xs), FromSlice(tc.ys)
		s1.DifferenceUpdate(s2)
		checkSet(t, s1, tc.want)
	}
	for _, tc := range tests {
		s1, s2 := FromSlice(tc.xs), FromSlice(tc.ys)
		checkSet(t, Difference(s1, s2), tc.want)
	}
}

func TestSymmetricDifference(t *testing.T) {
	xs := []int{3, 5, 11, 77}
	ys := []int{-1, 2, 6, 43}
	zs := append(ys, xs...)
	tests := []struct {
		xs, ys []int
		want   []int
	}{
		{nil, nil, nil},
		{[]int{}, []int{}, nil},
		{xs, xs, nil},
		// symmetric difference of disjoint is the same as the differnce
		{zs, xs, ys},
		{zs, ys, xs},
		//symmetric difference of disjoint sets is the union
		{xs, ys, zs},
		//
		{[]int{1, 3, 5, 9, 11}, []int{-2, 3, 9, 14, 22}, []int{-2, 1, 5, 11, 14, 22}},
	}
	for _, tc := range tests {
		s1, s2 := FromSlice(tc.xs), FromSlice(tc.ys)
		s1.SymmetricDifferenceUpdate(s2)
		if !checkSet(t, s1, tc.want) {
			t.Errorf("%v.xor.%v got = %v want=%v", tc.xs, tc.ys, s1.Slice(), tc.want)
		}
	}
	for _, tc := range tests {
		s1, s2 := FromSlice(tc.xs), FromSlice(tc.ys)
		if !checkSet(t, SymmetricDifference(s1, s2), tc.want) {
			t.Errorf("%v.xor.%v got = %v want=%v", tc.xs, tc.ys, s1.Slice(), tc.want)
		}
	}
}

func TestUnion(t *testing.T) {
	xs := []int{3, 5, 11, 77}
	ys := []int{-1, 2, 6, 43}
	zs := append(ys, xs...)
	tests := []struct {
		xs, ys []int
		want   []int
	}{
		{nil, nil, nil},
		{[]int{}, []int{}, nil},
		{xs, []int{}, xs},
		{xs, xs, xs}, // s | s = s
		{xs, ys, zs}, // union of disjoint set

		{zs, xs, zs}, // xs is a subset of zs
		{zs, ys, zs}, // ys is a subset of zs
	}
	for _, tc := range tests {
		s1, s2 := FromSlice(tc.xs), FromSlice(tc.ys)
		s1.Update(s2)
		if !checkSet(t, s1, tc.want) {
			t.Errorf("%v.union.%v got = %v want=%v", tc.xs, tc.ys, s1.Slice(), tc.want)
		}
		if !checkSet(t, Union(s1, s2), tc.want) {
			t.Errorf("%v.union.%v got = %v want=%v", tc.xs, tc.ys, s1.Slice(), tc.want)
		}
	}
}

func TestEquality(t *testing.T) {
	xs := FromSlice([]int{1, 2, 3})
	tests := []struct {
		xs, ys USet[int]
		want   bool
	}{
		{nil, nil, true},
		{xs, xs, true},
		{xs, xs.Copy(), true},
		{xs, nil, false},
		{xs, USet[int]{}, false},
		{USet[int]{}, USet[int]{}, true},
		{FromSlice([]int{1}), FromSlice([]int{2}), false},
		{FromSlice([]int{1, 3}), FromSlice([]int{2, 4}), false},
		{FromSlice([]int{1, 2, 3}), xs, true},
	}
	for i, tc := range tests {
		s1, s2 := tc.xs, tc.ys
		if g := s1.Equal(s2); g != tc.want {
			t.Errorf("case %d: %v.Equal(%v) got = %v want = %v", i, s1.Slice(), s2.Slice(), g, tc.want)
		}
	}
}

func TestForEach(t *testing.T) {
	xs := FromSlice([]int{1, 2, 3})
	sum := 0
	xs.ForEach(func(v int) { sum += v })
	if sum != 6 {
		t.Errorf("sum got %v want %v", sum, 6)
	}
	checkSet(t, xs, xs.Slice())
}
