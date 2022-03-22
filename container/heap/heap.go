package heap

import "sort"

type Interface interface {
	sort.Interface
	Push(x any) // add x as element Len()
	Pop() any   // remove and return element Len() - 1.
}

func Push(h Interface, x any) {
	h.Push(x)
	up(h, h.Len()-1)
}

func Pop(h Interface) any {
	n := h.Len() - 1
	h.Swap(0, n)
	down(h, 0, n)
	return h.Pop()
}

func Init(h Interface) {
	// heapify
	n := h.Len()
	for i := 0; i < n; i++ {
		down(h, i, n)
	}
}

func Remove(h Interface, i int) any {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !down(h, i, n) {
			up(h, i)
		}
	}
	return h.Pop()
}

func Fix(h Interface, i int) {
	if !down(h, i, h.Len()) {
		up(h, i)
	}
}

func up(h Interface, j int) {
	for {
		i := (j - 1) / 2 //parent
		if i == j || h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}
func down(h Interface, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}
