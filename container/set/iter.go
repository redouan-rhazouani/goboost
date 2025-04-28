package set

import (
	"iter"
	"maps"
)

// All returns an iterator over elements from s.
// The iteration order is not specified and is not guaranteed
// to be the same from one call to the next.
func (s HashSet[T]) All() iter.Seq[T] {
	return maps.Keys(s.m)
}

// Insert adds the elements from seq to s.
// If a key in seq already exists in m, its value will be overwritten.
func Insert[T comparable](s HashSet[T], seq iter.Seq[T]) {
	for v := range seq {
		s.m[v] = struct{}{}
	}
}

// Union visits the values representing the union
//
// that is, the values in s1 or s2, without duplicates.
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

// Intersection visits the values representing the intersection
//
//	i.e., the values that are both in s1 and s2.
func Intersection[T comparable](s1, s2 HashSet[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		s1, s2 = swapIfLess(s1, s2)
		for v := range s1.m {
			if _, ok := s2.m[v]; ok {
				if !yield(v) {
					return
				}
			}
		}
	}
}

//	Difference Iterates over the values in the difference between s1 and s2
//
// i.e., that is, values present in either s1 or s2, but not in both
func Difference[T comparable](s1, s2 HashSet[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s1.m {
			if _, ok := s2.m[k]; !ok {
				if !yield(k) {
					return
				}
			}
		}
	}
}

// SymmetricDifference Iterates over the values in the symmetric difference between s1 and s2
//
// that is, values present in either s1 or s2, but not in both
func SymmetricDifference[T comparable](s, o HashSet[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.m {
			if _, ok := o.m[k]; !ok {
				if !yield(k) {
					return
				}
			}
		}
		for k := range o.m {
			if _, ok := s.m[k]; !ok {
				if !yield(k) {
					return
				}
			}
		}
	}
}

// Collect collects elements from seq into a new HashSet.
func Collect[T comparable](seq iter.Seq[T]) HashSet[T] {
	s := Make[T]()
	Insert(s, seq)
	return s
}
