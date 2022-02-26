// Package vector implements dynamic size arrays
package vector

// Vector is a sequence container that encapsulates dynamic size arrays
type Vector[T any] []T

// Make construct a vector with specified length and capacity.
//
//	The capacity must be no smaller than the
//	length. For example, Make[T](2, 10) allocates an underlying array
//	of size 10 and returns a vector of length 2 and capacity 10 that is
//	backed by this underlying array.
func Make[T any](length, capacity int) Vector[T] {
	return make([]T, length, capacity)
}

// Capacity returns the number of elements the vector can hold without allocating
func (vec Vector[T]) Capacity() int {
	return cap(vec)
}

// Len returns the number of elements int the vector
func (vec Vector[T]) Len() int {
	return len(vec)
}

// At returns element at position i
// i < 0 means access element at len(vec) - 1 (At(-1) is the last element)
func (vec Vector[T]) At(i int) T {
	if i < 0 {
		i = len(vec) + i
	}
	return vec[i]
}

// IndexFunc returns the index into vec of the first element
// satisfying f(x), or -1 if none do.
func (vec Vector[T]) IndexFunc(f func(T) bool) int {
	for i, x := range vec {
		if f(x) {
			return i
		}
	}
	return -1
}

// LastIndexFunc returns the index into vec of the last occurence of element
// satisfying f(x), or -1 if none do.
func (vec Vector[T]) LastIndexFunc(f func(T) bool) int {
	for i := vec.Len() - 1; i >= 0; i-- {
		if f(vec[i]) {
			return i
		}
	}
	return -1
}

// Index returns the index of the first occurrence of v in s, or -1 if substr is not present in s.
func Index[S ~[]T, T comparable](s S, v T) int {
	for i, x := range s {
		if v == x {
			return i
		}
	}
	return -1
}

// LastIndex returns the index of the last occurrence of v in s, or -1 if substr is not present in s.
func LastIndex[S ~[]T, T comparable](s S, v T) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == v {
			return i
		}
	}
	return -1
}

func Equal[S ~[]T, T comparable](x, y S) bool {
	if len(x) != len(y) {
		return false
	}
	for i, x := range x {
		if x != y[i] {
			return false
		}
	}
	return true
}

// Copy returns a shallow copy of vec
func (vec Vector[T]) Copy() Vector[T] {
	buf := make(Vector[T], len(vec))
	copy(buf, vec)
	return buf
}

// Append appends elements to the end of a vector.
//
// Note: Append is just a wrapper around the builtin function append
func (vec *Vector[T]) Append(xs ...T) {
	*vec = append(*vec, xs...)
}

// Push appends x to the back of vec
func (vec *Vector[T]) Push(x T) {
	*vec = append(*vec, x)
}

func (vec *Vector[T]) Pop() T {
	s := *vec
	n := len(s) - 1
	var elem T
	s[n], elem = elem, s[n] //avoid memory leaks
	*vec = s[:n]
	return elem
}

// SwapDelete removes and returns the element at position i from vec.
//
// The removed element is replaced by the last element of the vector
// Operation doesn't preserve ordering, but is O(1).
// If you need to preserve ordering, use remove instead
// Panics if index is out of bounds
func (vec *Vector[T]) SwapDelete(i int) T {
	s := *vec
	elem, len := s[i], s.Len()
	s[i], s[len-1] = s[len-1], zero[T]() //avoid memory leaks
	*vec = s[:len-1]
	return elem
}

// Delete removes and returns element at position i within vec, shifting all
// elements after it to the left.
// Note: Because this shifts over the remaining elements, it has a worst-case performance of O(n).
// If you don’t need the order of elements to be preserved, use SwapRemove instead.
func (vec *Vector[T]) Delete(i int) T {
	s := *vec
	elem := s[i]
	*vec = append(s[:i], s[i+1:]...)
	s[len(s)-1] = zero[T]() //avoid memory leaks
	return elem
}

func (vec *Vector[T]) DeleteRange(low, high int) {
	s := *vec
	*vec = append(s[:low], s[high:]...)
	zeroRange(s, len(s)-(high-low), len(s)) //avoid memory leaks

}

// Clear clears all elements from vec.
// After this call the len() returns zero while the capacity stays unchanged.
// Note that this method set all elements to T{} to avoid memory leaks
func (vec *Vector[T]) Clear() {
	s := *vec
	zeroRange(s, 0, len(s))
	*vec = s[:0]
}

func (vec *Vector[T]) RemoveFunc(f func(v T) bool) {
	s := *vec
	i := s.IndexFunc(f)
	if i < 0 {
		return
	}
	j := i + 1
	for ; j < s.Len(); j = j + 1 {
		if !f(s[j]) {
			s[i] = s[j]
			i++
		}
	}
	*vec = s[:i]
	zeroRange(s, i, j)
}

func Remove[S ~[]T, T comparable](p *S, v T) {
	s := *p
	i := Index(s, v)
	if i < 0 {
		return
	}
	j := i + 1
	for ; j < len(s); j = j+1 {
		if s[j] != v {
			s[i] = s[j]
			i++
		}
	}
	zeroRange[T](s, i, j)
	*p = s[:i]
}

func Count[S ~[]T, T comparable](s S, v T) int {
	n := 0
	for _, x := range s {
		if x == v {
			n++
		}
	}
	return n
}

// zero returns T's zero value (T{})
func zero[T any]() T { var zero T; return zero }

func zeroRange[T any](s []T, low, high int) {
	var zero T
	for ; low < high; low++ {
		s[low] = zero //avoid memory leaks
	}
}
