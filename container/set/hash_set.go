// Package set implements hash sets of any comparable type
package set

// Set is an unordered collection with unique elements.
//
// Set requires that the elements statisfy the comparable constraint.
// It support mathematical operations like union, difference, and symmetric differce
type Set[T comparable] map[T]struct{}

// Make creates an empty set
//
// The empty set is allocated with enough space to hold the
// specified number of elements.
func Make[T comparable]() Set[T] {
	return make(Set[T])
}

// Make creates an empty set with the specified capacity cap
func MakeWithCapacity[T comparable](cap int) Set[T] {
	return make(Set[T], cap)
}

// FromSlice creates a new set using elements of xs
func FromSlice[T comparable](xs []T) Set[T] {
	s := MakeWithCapacity[T](len(xs))
	for _, v := range xs {
		s[v] = struct{}{}
	}
	return s
}

// Len returns the number of elements of set s
func (s Set[T]) Len() int {
	return len(s)
}

// Add element v to the set s
// if v is in s this has no effect
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

// InsertSlice inserts elements from slice xs
func (s Set[T]) InsertSlice(xs []T) {
	for _, v := range xs {
		s[v] = struct{}{}
	}
}

// Delete element v from the set s
// If v not in s this has no effect
func (s Set[T]) Delete(v T) {
	delete(s, v)
}

// DeleteIF deletes all elements that satisfy the predicate pred from set
func (s Set[T]) DeleteIF(pred func(T) bool) {
	for v := range s {
		if pred(v) {
			delete(s, v)
		}
	}
}

// Contains reports whether v is in s
func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// Do applies function f to each element of s
// It's ok for to call s.Delete(v)
func (s Set[T]) Do(f func(T)) {
	for v := range s {
		f(v)
	}
}

// Update set s, adding elements from set o
func (s Set[T]) Update(o Set[T]) {
	for v := range o {
		s[v] = struct{}{}
	}
}

// IntersectionUpdate set, keeping only common elements between s and o
func (s Set[T]) IntersectionUpdate(o Set[T]) {
	for v := range s {
		if _, ok := o[v]; !ok {
			delete(s, v)
		}
	}
}

// DifferenceUpdate the set, removing elements found in o
func (s Set[T]) DifferenceUpdate(o Set[T]) {
	s1, s2 := swapIfLess(s, o)
	for v := range s1 {
		if _, ok := s2[v]; ok {
			delete(s, v)
		}
	}
}

// SymmetricDifferenceUpdate, keeping elements found in either s or o but not in both
func (s Set[T]) SymmetricDifferenceUpdate(o Set[T]) {
	complement := make(Set[T])
	for v := range s {
		if _, ok := o[v]; ok {
			complement[v] = struct{}{}
			delete(s, v)
		}
	}
	for v := range o {
		if _, ok := complement[v]; !ok {
			s.Add(v)
		}
	}
}

// IsDisjoint return true if sets s and o has no element in common.
// Two sets are disjoint if and only if their intersection is the empty set
func (s Set[T]) IsDisjoint(o Set[T]) bool {
	s1, s2 := swapIfLess(s, o)
	for v := range s1 {
		if _, ok := s2[v]; ok {
			return false
		}
	}
	return true
}

// IsSubset test whether every element in s is also in o
func (s Set[T]) IsSubset(o Set[T]) bool {
	if s.Len() > o.Len() {
		return false
	}
	for v := range s {
		if _, ok := o[v]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset test whether every element in o is also in s
func (s Set[T]) IsSuperset(o Set[T]) bool {
	return o.IsSubset(s)
}

// Copy return a copy of the set
func (s Set[T]) Copy() Set[T] {
	c := make(Set[T], len(s))
	for v := range s {
		c[v] = struct{}{}
	}
	return c
}

// Clear removes all elements from the set
func (s Set[T]) Clear() {
	for v := range s {
		s.Delete(v)
	}
}

// Equal returns true if the contents are equal
func (s Set[T]) Equal(o Set[T]) bool {
	if len(s) != len(o) {
		return false
	}
	for v := range s {
		if _, ok := o[v]; !ok {
			return false
		}
	}
	return true
}

// Slice returns set elements as a slice
func (s Set[T]) Slice() []T {
	a := make([]T, len(s))
	i := 0
	for v := range s {
		a[i] = v
		i++
	}
	return a
}

// Union returns a new set with elements from set s and o
func Union[T comparable](s, o Set[T]) Set[T] {
	u := make(Set[T])
	for v := range s {
		u[v] = struct{}{}
	}
	for v := range o {
		u[v] = struct{}{}
	}
	return u
}

// Intersection returns new set with elements common to set s and o
func Intersection[T comparable](s, o Set[T]) Set[T] {
	intersect := make(Set[T])
	s, o = swapIfLess(s, o)
	for v := range s {
		if _, ok := o[v]; ok {
			intersect[v] = struct{}{}
		}
	}
	return intersect
}

// swapIfLess swap both sets if len(s) < len(o)
func swapIfLess[T comparable](s, o Set[T]) (Set[T], Set[T]) {
	if len(s) < len(o) {
		return s, o
	}
	return o, s
}

// Difference returns new set with elements in the set s that are not in o
func Difference[T comparable](s, o Set[T]) Set[T] {
	diff := make(Set[T])
	for v := range s {
		if _, ok := o[v]; !ok {
			diff[v] = struct{}{}
		}
	}
	return diff
}

// SymmetricDifference returns new set with elements in either s or o but not both
func SymmetricDifference[T comparable](s, o Set[T]) Set[T] {
	diff := make(Set[T])
	for v := range s {
		if _, ok := o[v]; !ok {
			diff[v] = struct{}{}
		}
	}
	for v := range o {
		if _, ok := s[v]; !ok {
			diff[v] = struct{}{}
		}
	}
	return diff
}
