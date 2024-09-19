package vector

// Partition rearranges the elements of the sequence 'a' based on the predicate 'f'.
// Elements for which 'f' returns true are moved to the front, and elements for
// which it returns false are moved to the back. The relative order of the elements
// within each partition is not preserved.
// It returns the index where the partition is split, i.e., the first element in the
// sequence that evaluates to false
func Partition[S ~[]E, E any](a S, f func(i int) bool) int {
	first := 0
	for ; first < len(a); first++ {
		if !f(first) {
			break
		}
	}
	if first == len(a) {
		return first
	}

	for i := first + 1; i < len(a); i++ {
		if f(i) {
			a[i], a[first] = a[first], a[i]
			first++
		}
	}
	return first
}

// IsPartitioned checks if the sequence [0, n) is partitioned
// according to the predicate 'f'. It returns true if all elements that satisfy
// 'f' precede those that do not, otherwise false.
func IsPartitioned(n int, f func(i int) bool) bool {
	first := 0
	for ; first < n; first++ {
		if !f(first) {
			break
		}
	}

	for ; first < n; first++ {
		if f(first) {
			return false
		}
	}
	return true
}

// Examines the partitioned range [0, n) and returns the index of the first element
// that does not satisfy the predicate 'p', or 'n' if all elements satisfy 'p'.
// If the elements in the range [0, n) are not partitioned according to the predicate
// 'f' the behavior is undefined.
func PartitionPoint(n int, f func(i int) bool) int {
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		if f(h) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
