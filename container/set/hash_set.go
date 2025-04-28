// Package set implements hash sets of any comparable type
package set

import (
	"maps"
)

// HashSet is an unordered collection with unique elements.
//
// HashSet requires that the elements satisfy the comparable constraint.
// It support mathematical operations like union, difference, and symmetric difference
type HashSet[T comparable] struct {
	m map[T]struct{}
}

// Make creates an empty set
//
// The empty set is allocated with enough space to hold the
// specified number of elements.
func Make[T comparable]() HashSet[T] {
	return HashSet[T]{
		m: make(map[T]struct{}),
	}
}

// Make creates an empty set with the specified capacity cap
func MakeWithCapacity[T comparable](cap int) HashSet[T] {
	return HashSet[T]{
		m: make(map[T]struct{}, cap),
	}
}

// Len returns the number of elements of set s
func (s HashSet[T]) Len() int {
	return len(s.m)
}

// Copy return a copy of the set
func (s HashSet[T]) Clone() HashSet[T] {
	return HashSet[T]{
		m: maps.Clone(s.m),
	}
}

// Clear removes all elements from the set
func (s HashSet[T]) Clear() {
	clear(s.m)
}

// Contains reports whether v is in s
func (s HashSet[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// FromSlice creates a new set using elements of xs
func FromSlice[T comparable](xs []T) HashSet[T] {
	s := MakeWithCapacity[T](len(xs))
	for _, v := range xs {
		s.m[v] = struct{}{}
	}
	return s
}

// Add element v to the set s
// if v is in s this has no effect
func (s HashSet[T]) Add(v T) {
	s.m[v] = struct{}{}
}

// Delete element v from the set s
// If v not in s this has no effect
func (s HashSet[T]) Delete(v T) {
	delete(s.m, v)
}

// DeleteFunc deletes all elements that satisfy the predicate pred from set
func (s HashSet[T]) DeleteFunc(pred func(T) bool) {
	maps.DeleteFunc(s.m, func(k T, _ struct{}) bool { return pred(k) })
}

// Copy copies all elements in src adding them to dst.
func Copy[E comparable](dst, src HashSet[E]) {
	maps.Copy(dst.m, src.m)
}

// IntersectionUpdate set, keeping only common elements between s and other
func (s HashSet[T]) IntersectionUpdate(other HashSet[T]) {
	for v := range s.m {
		if _, ok := other.m[v]; !ok {
			delete(s.m, v)
		}
	}
}

// DifferenceUpdate the set, removing elements found in other
func (s HashSet[T]) DifferenceUpdate(other HashSet[T]) {
	s1, s2 := swapIfLess(s, other)
	for v := range s1.m {
		if _, ok := s2.m[v]; ok {
			delete(s.m, v)
		}
	}
}

// SymmetricDifferenceUpdate, keeping elements found in either s or other
// , that is, values present in either s or other,but not in both
func (s HashSet[T]) SymmetricDifferenceUpdate(other HashSet[T]) {
	complement := make(map[T]struct{})
	for v := range s.m {
		if _, ok := other.m[v]; ok {
			complement[v] = struct{}{}
			delete(s.m, v)
		}
	}
	for v := range other.m {
		if _, ok := complement[v]; !ok {
			s.Add(v)
		}
	}
}

// IsDisjoint return true if sets s and other has no element in common.
// Two sets are disjoint if and only if their intersection is the empty set
func (s HashSet[T]) IsDisjoint(other HashSet[T]) bool {
	s1, s2 := swapIfLess(s, other)
	for v := range s1.m {
		if _, ok := s2.m[v]; ok {
			return false
		}
	}
	return true
}

// IsSubset test whether every element in s is also in other
func (s HashSet[T]) IsSubset(other HashSet[T]) bool {
	if s.Len() > other.Len() {
		return false
	}
	for v := range s.m {
		if _, ok := other.m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset test whether every element in other is also in s
func (s HashSet[T]) IsSuperset(other HashSet[T]) bool {
	return other.IsSubset(s)
}

// Equal returns true if the contents are equal
func (s HashSet[T]) Equal(other HashSet[T]) bool {
	return maps.Equal(s.m, other.m)
}

// Slice returns set elements as a slice
func (s HashSet[T]) Slice() []T {
	a := make([]T, len(s.m))
	i := 0
	for v := range s.m {
		a[i] = v
		i++
	}
	return a
}

// swapIfLess swap both sets if len(s) < len(o)
func swapIfLess[T comparable](s, o HashSet[T]) (HashSet[T], HashSet[T]) {
	if len(s.m) < len(o.m) {
		return s, o
	}
	return o, s
}
