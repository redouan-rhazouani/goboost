package set

import (
	"iter"
	"maps"
)

// All returns an iterator over keys from s.
// The iteration order is not specified and is not guaranteed
// to be the same from one call to the next.
func (s HashSet[T]) All() iter.Seq[T] {
	return maps.Keys(s.m)
}

// Insert adds the elements from seq to m.
// If a key in seq already exists in m, its value will be overwritten.
func Insert[T comparable](s HashSet[T], seq iter.Seq[T]) {
	for v := range seq {
		s.m[v] = struct{}{}
	}
}

// Union returns a new set with elements from set s and o
func Union[T comparable](s, o HashSet[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		s, o = swapIfLess(s, o)
		for v := range o.m {
			if !yield(v) {
				return
			}
		}

		for v := range s.m {
			if _, ok := o.m[v]; ok {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Intersection returns new set with elements common to set s and o
func Intersection[T comparable](s, o HashSet[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		s, o = swapIfLess(s, o)
		for v := range s.m {
			if _, ok := o.m[v]; ok {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Difference returns new set with elements in the set s that are not in o
func Difference[T comparable](s, o HashSet[T]) iter.Seq[T] {
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
func SymmetricDifference[T comparable](s, o HashSet[T]) iter.Seq[T] {
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

// Collect collects elements from seq into a new set and returns it.
func Collect[T comparable](seq iter.Seq[T]) HashSet[T] {
	s := Make[T]()
	Insert(s, seq)
	return s
}
