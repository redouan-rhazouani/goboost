package vector

func Partition(a []int, pred func(i int) bool) int {
	i, j := 0, len(a)-1
	for i < j {
		for i < len(a) && pred(i) {
			i++
		}
		for j > i && !pred(j) {
			j--
		}

		if i < j {
			a[i], a[j] = a[j], a[i]
			j--
			i++
		}

	}
	return i
}

// Examines the partitioned range [0, n) and locates the end of the first partition, that is, the first element that does not satisfy p or last if all elements satisfy p.
// If the elements elem of [0, n) are not partitioned with respect to the expression bool(p(elem)), the behavior is undefined.
func PartitionPoint(a []int, pred func(i int) bool) int {
	length := len(a)
	first := 0
	for length > 0 {
		half := length / 2
		mid := first + half
		if pred(mid) {
			first = mid + 1
			length -= (half + 1)
		} else {
			length = half
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
