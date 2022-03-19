package priority_queue

import "constraints"

type Compare[T any] func(T, T) bool

func Less[T constraints.Ordered](a, b T) bool {
	return a < b
}

type PriorityQueue[T any] struct {
	s   []T
	cmp Compare[T]
}

func Make[T any](capacity int, cmp Compare[T]) PriorityQueue[T] {
	pq := PriorityQueue[T]{}
	pq.s = make([]T, 0, capacity)
	pq.cmp = cmp
	return pq
}

func (pq *PriorityQueue[T]) Len() int {
	return len(pq.s)
}

func (pq *PriorityQueue[T]) Push(x T) {
	pq.s = append(pq.s, x)
	siftUp(pq.s, len(pq.s)-1, pq.cmp)
}

func (pq *PriorityQueue[T]) Pop() T {
	s := pq.s
	n := len(s) - 1
	sn := s[n]
	s[0] = s[n]
	var zero T
	s[n] = zero // avoid memory leak
	siftDown(s, 0, n, pq.cmp)
	pq.s = s[:n]
	return sn
}

// Top accesses the top element
func (pq *PriorityQueue[T]) Top() T {
	return pq.s[pq.Len()-1]
}

func siftUp[T any](s []T, j int, cmp Compare[T]) {
	for {
		i := (j - 1) / 2 // parent
		if j == i || !cmp(s[j], s[i]) {
			break
		}
		s[j], s[i] = s[i], s[j]
		j = i
	}
}

func siftDown[T any](s []T, i0, n int, cmp Compare[T]) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int workflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && j2 >= 0 && cmp(s[j2], s[j1]) {
			j = j2 // right child
		}
		if !cmp(s[j], s[i]) {
			break
		}
		s[i], s[j] = s[j], s[i]
		i = j
	}
	return i > i0
}
