// Package sets implements sets of any comparable type
package sets

// USet is an unordered collection with unique elements.
// USet support mathematical operations like union, difference, and symmetric differce
type USet[T comparable] map[T]struct{}

// Make creates an empty set
//
// The empty set is allocated with enough space to hold the
// specified number of elements.
func Make[T comparable]() USet[T] {
	return make(USet[T])
}

// Make creates an empty set with the specified capacity cap
func MakeWithCapacity[T comparable](cap int) USet[T] {
	return make(USet[T], cap)
}

// FromSlice creates a new set using elements from xs
func FromSlice[T comparable](xs []T) USet[T] {
	s := MakeWithCapacity[T](len(xs))
	for _, v := range xs {
		s[v] = struct{}{}
	}
	return s
}

// Len returns the number of elements of set s
func (s USet[T]) Len() int {
	return len(s)
}

// Add element v to the set s
// if v is in s this has no effect
func (s USet[T]) Add(v T) {
	s[v] = struct{}{}
}

// InsertSlice inserts elements from slice xs
func (s USet[T]) InsertSlice(xs []T) {
	for _, v := range xs {
		s[v] = struct{}{}
	}
}

// Delete element v from the set s
// If v not in s this has no effect
func (s USet[T]) Delete(v T) {
	delete(s, v)
}

// DeleteIF deletes all elements that satisfy the predicate pred from set
func (s USet[T]) DeleteIF(pred func(T) bool) {
	for v := range s {
		if pred(v) {
			delete(s, v)
		}
	}
}

// Contains reports whether v is in s
func (s USet[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// ForEach applies function f to each element of s
// It's ok for to call s.Delete(v)
func (s USet[T]) ForEach(f func(T)) {
	for v := range s {
		f(v)
	}
}

// Update set s, adding elements from set o
func (s USet[T]) Update(o USet[T]) {
	for v := range o {
		s[v] = struct{}{}
	}
}

// IntersectionUpdate set, keeping only common elements between s and o
func (s USet[T]) IntersectionUpdate(o USet[T]) {
	for v := range s {
		if _, ok := o[v]; !ok {
			delete(s, v)
		}
	}
}

// DifferenceUpdate the set, removing elements that are not in o
func (s USet[T]) DifferenceUpdate(o USet[T]) {
	for v := range s {
		if _, ok := o[v]; ok {
			delete(s, v)
		}
	}
}

// SymmetricDifferenceUpdate, keeping elements found in either s or o but not in both
func (s USet[T]) SymmetricDifferenceUpdate(o USet[T]) {
	complement := make(USet[T])
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
// Two sets are disjoint if and only their intersection is the empty set
func (s USet[T]) IsDisjoint(o USet[T]) bool {
	for v := range s {
		if _, ok := o[v]; ok {
			return false
		}
	}
	return true
}

// IsSubset test whether every element is s is also in o
func (s USet[T]) IsSubset(o USet[T]) bool {
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

// IsSuperset test whether every element is o is also in s
func (s USet[T]) IsSuperset(o USet[T]) bool {
	return o.IsSubset(s)
}

// Copy return a copy of the set
func (s USet[T]) Copy() USet[T] {
	c := make(USet[T], len(s))
	for v := range s {
		c[v] = struct{}{}
	}
	return c
}

// Clear removes all elements from the set
func (s USet[T]) Clear() {
	for v := range s {
		s.Delete(v)
	}
}

// Equal returns true if the contents are equal
func (s USet[T]) Equal(o USet[T]) bool {
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

func (s USet[T]) Slice() []T {
	a := make([]T, len(s))
	i := 0
	for v := range s {
		a[i] = v
		i++
	}
	return a
}

// Union returns a new set with elements from set s and o
func Union[T comparable](s, o USet[T]) USet[T] {
	u := make(USet[T])
	for v := range s {
		u[v] = struct{}{}
	}
	for v := range o {
		u[v] = struct{}{}
	}
	return u
}

// Intersection returns new set with elements common to set s and o
func Intersection[T comparable](s, o USet[T]) USet[T] {
	intersect := make(USet[T])
	for v := range s {
		if _, ok := o[v]; ok {
			intersect[v] = struct{}{}
		}
	}
	return intersect
}

// Difference returns new set with elements in the set s that are not in o
func Difference[T comparable](s, o USet[T]) USet[T] {
	diff := make(USet[T])
	for v := range s {
		if _, ok := o[v]; !ok {
			diff[v] = struct{}{}
		}
	}
	return diff
}

// SymmetricDifference returns new set with elements in either s or o but not both
func SymmetricDifference[T comparable](s, o USet[T]) USet[T] {
	diff := make(USet[T])
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
