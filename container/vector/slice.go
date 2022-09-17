package vector

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

// Count returns the number of occurence of v in s
func Count[S ~[]T, T comparable](s S, v T) int {
	n := 0
	for _, x := range s {
		if x == v {
			n++
		}
	}
	return n
}

func Remove[S ~[]T, T comparable](p *S, v T) {
	s := *p
	i := Index(s, v)
	if i < 0 {
		return
	}
	j := i + 1
	for ; j < len(s); j = j + 1 {
		if s[j] != v {
			s[i] = s[j]
			i++
		}
	}
	zeroRange(s, i, j)
	*p = s[:i]
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
