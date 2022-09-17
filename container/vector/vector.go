// Package vector implements a contiguous growable array type
package vector

// Vector is a wrapper around a generic slice.
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

// Cap returns the number of elements the vector can hold without allocating
func (vec Vector[T]) Cap() int {
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

// RemoveFunc removes all elements satisfying f()
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

// zero returns T's zero value (T{})
func zero[T any]() T { var zero T; return zero }

func zeroRange[T any](s []T, low, high int) {
	var zero T
	for ; low < high; low++ {
		s[low] = zero //avoid memory leaks
	}
}
