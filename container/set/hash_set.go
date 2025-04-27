// Package set implements hash sets of any comparable type
package set

import (
	"iter"
	"maps"
)

// Set is an unordered collection with unique elements.
//
// Set requires that the elements satisfy the comparable constraint.
// It support mathematical operations like union, difference, and symmetric difference
type Set[T comparable] struct {
	m map[T]struct{}
}

// Make creates an empty set
//
// The empty set is allocated with enough space to hold the
// specified number of elements.
func Make[T comparable]() Set[T] {
	return Set[T]{
		m: make(map[T]struct{}),
	}
}

// Make creates an empty set with the specified capacity cap
func MakeWithCapacity[T comparable](cap int) Set[T] {
	return Set[T]{
		m: make(map[T]struct{}, cap),
	}
}

// Len returns the number of elements of set s
func (s Set[T]) Len() int {
	return len(s.m)
}

// All returns an iterator over keys from s.
// The iteration order is not specified and is not guaranteed
// to be the same from one call to the next.
func (s Set[T]) All() iter.Seq[T] {
	return maps.Keys(s.m)
}

// Copy return a copy of the set
func (s Set[T]) Clone() Set[T] {
	return Set[T]{
		m: maps.Clone(s.m),
	}
}

// Clear removes all elements from the set
func (s Set[T]) Clear() {
	clear(s.m)
}

// FromSlice creates a new set using elements of xs
func FromSlice[T comparable](xs []T) Set[T] {
	s := MakeWithCapacity[T](len(xs))
	for _, v := range xs {
		s.m[v] = struct{}{}
	}
	return s
}

// Add element v to the set s
// if v is in s this has no effect
func (s Set[T]) Add(v T) {
	s.m[v] = struct{}{}
}

// Insert adds the elements from seq to m.
// If a key in seq already exists in m, its value will be overwritten.
func (s Set[T]) Insert(seq iter.Seq[T]) {
	for v := range seq {
		s.m[v] = struct{}{}
	}
}

// Delete element v from the set s
// If v not in s this has no effect
func (s Set[T]) Delete(v T) {
	delete(s.m, v)
}

// DeleteFunc deletes all elements that satisfy the predicate pred from set
func (s Set[T]) DeleteFunc(pred func(T) bool) {
	maps.DeleteFunc(s.m, func(k T, _ struct{}) bool { return pred(k) })
}

// Contains reports whether v is in s
func (s Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Update set s, adding elements from set o
func (s Set[T]) Update(o Set[T]) {
	for v := range o.m {
		s.m[v] = struct{}{}
	}
}

// IntersectionUpdate set, keeping only common elements between s and o
func (s Set[T]) IntersectionUpdate(o Set[T]) {
	for v := range s.m {
		if _, ok := o.m[v]; !ok {
			delete(s.m, v)
		}
	}
}

// DifferenceUpdate the set, removing elements found in o
func (s Set[T]) DifferenceUpdate(o Set[T]) {
	s1, s2 := swapIfLess(s, o)
	for v := range s1.m {
		if _, ok := s2.m[v]; ok {
			delete(s.m, v)
		}
	}
}

// SymmetricDifferenceUpdate, keeping elements found in either s or o but not in both
func (s Set[T]) SymmetricDifferenceUpdate(o Set[T]) {
	complement := make(map[T]struct{})
	for v := range s.m {
		if _, ok := o.m[v]; ok {
			complement[v] = struct{}{}
			delete(s.m, v)
		}
	}
	for v := range o.m {
		if _, ok := complement[v]; !ok {
			s.Add(v)
		}
	}
}

// IsDisjoint return true if sets s and o has no element in common.
// Two sets are disjoint if and only if their intersection is the empty set
func (s Set[T]) IsDisjoint(o Set[T]) bool {
	s1, s2 := swapIfLess(s, o)
	for v := range s1.m {
		if _, ok := s2.m[v]; ok {
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
	for v := range s.m {
		if _, ok := o.m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset test whether every element in o is also in s
func (s Set[T]) IsSuperset(o Set[T]) bool {
	return o.IsSubset(s)
}

// Equal returns true if the contents are equal
func (s Set[T]) Equal(o Set[T]) bool {
	if len(s.m) != len(o.m) {
		return false
	}
	for v := range s.m {
		if _, ok := o.m[v]; !ok {
			return false
		}
	}
	return true
}

// Slice returns set elements as a slice
func (s Set[T]) Slice() []T {
	a := make([]T, len(s.m))
	i := 0
	for v := range s.m {
		a[i] = v
		i++
	}
	return a
}

// Union returns a new set with elements from set s and o
func Union[T comparable](s, o Set[T]) Set[T] {
	u := make(map[T]struct{})
	for v := range s.m {
		u[v] = struct{}{}
	}
	for v := range o.m {
		u[v] = struct{}{}
	}
	return Set[T]{m: u}
}

// Intersection returns new set with elements common to set s and o
func Intersection[T comparable](s, o Set[T]) Set[T] {
	intersect := Make[T]()
	s, o = swapIfLess(s, o)
	for v := range s.m {
		if _, ok := o.m[v]; ok {
			intersect.m[v] = struct{}{}
		}
	}
	return intersect
}

// swapIfLess swap both sets if len(s) < len(o)
func swapIfLess[T comparable](s, o Set[T]) (Set[T], Set[T]) {
	if len(s.m) < len(o.m) {
		return s, o
	}
	return o, s
}

// Difference returns new set with elements in the set s that are not in o
func Difference[T comparable](s, o Set[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s.m {
			if _, ok := o.m[v]; !ok {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// SymmetricDifference returns new set with elements in either s or o but not both
func SymmetricDifference[T comparable](s, o Set[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s.m {
			if _, ok := o.m[v]; !ok {
				if !yield(v) {
					return
				}
			}
		}
		for v := range o.m {
			if _, ok := s.m[v]; !ok {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Construct[T comparable](seq iter.Seq[T]) Set[T] {
	s := Make[T]()
	for v := range seq {
		s.Add(v)
	}
	return s
}
