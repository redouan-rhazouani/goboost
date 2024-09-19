package vector

import "slices"

// LastIndex returns the index of the last occurrence of v in s, or -1 if v is not present in s.
func LastIndex[S ~[]T, T comparable](s S, v T) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// Count returns the number of occurrence of v in s
func Count[S ~[]T, T comparable](s S, v T) int {
	n := 0
	for _, x := range s {
		if x == v {
			n++
		}
	}
	return n
}

func Remove[S ~[]E, E comparable](p S, v E) S {
	return slices.DeleteFunc(p, func(x E) bool {
		return x == v
	})
}
