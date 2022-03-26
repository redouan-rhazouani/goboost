package forward_list

import (
	"testing"
)

func checkListLen[T any](t *testing.T, l *ForwardList[T], len int) bool {
	t.Helper()
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers[T any](t *testing.T, l *ForwardList[T], es []*Element[T]) {
	root := &l.root
	t.Helper()
	if !checkListLen(t, l, len(es)) {
		return
	}

	// zero length lists must be the zero value or properly initialized
	if len(es) == 0 {
		if l.root.next != nil && l.root.next != root || l.back != root {
			t.Errorf("l.root.next = %p, l.back = %p; both should both be nil or %p", l.root.next, l.back, root)
		}
		return
	}

	// check internal and external connections
	for i, e := range es {
		next := root // l.back.next = root
		Next := (*Element[T])(nil)
		if i < len(es)-1 {
			next = es[i+1]
			Next = next
		}
		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

func checkList[T comparable](t *testing.T, l *ForwardList[T], es []T) {
	t.Helper()
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, e.Value, es[i])
		}
		i++
	}
}

func TestList(t *testing.T) {
	{
		// empty list
		l := New[string]()
		checkListPointers(t, l, []*Element[string]{})
		if l.Back() != l.Front() || l.Back() != nil {
			t.Errorf("Empty list back and front must point to nil")
		}

		// Single Element[T] list
		e := l.PushFront("a")
		checkListPointers(t, l, []*Element[string]{e})
	}
	l := New[int]()
	// Bigger list
	e2 := l.PushFront(2)
	e1 := l.PushFront(1)
	e3 := l.PushBack(3)
	e4 := l.PushBack(0)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})

	l.RemoveAfter(e1)
	checkListPointers(t, l, []*Element[int]{e1, e3, e4})

	e2 = l.InsertAfter(2, e1) // insert after front
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})
	l.RemoveAfter(e1)
	e2 = l.InsertAfter(2, e3) // insert after middle
	checkListPointers(t, l, []*Element[int]{e1, e3, e2, e4})
	l.RemoveAfter(e3)
	e2 = l.InsertAfter(2, e4) // insert after back
	checkListPointers(t, l, []*Element[int]{e1, e3, e4, e2})
	l.RemoveAfter(e4)

	// Check standard iteration.
	sum := 0
	for e := l.Front(); e != nil; e = e.Next() {
		sum += e.Value
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// // Clear all Element[T]s by iterating
	// for e := l.Front(); e != nil; {
	// 	e = l.Remove(e)
	// }
	// checkListPointers(t, l, []*Element[int]{})
	// e1 = l.PushBack(1)
	// e2 = l.PushBack(2)
	// checkListPointers(t, l, []*Element[int]{e1, e2})
	// for e := l.Back(); e != nil; {
	// 	e = l.Remove(e)
	// }
	// checkListPointers(t, l, []*Element[int]{})
}

func TestExtending(t *testing.T) {
	l1 := New[int]()
	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2 := New[int]()
	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New[int]()
	l3.PushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushBackList(l2)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	l3.Clear()
	l3.PushFrontList(l2)
	checkList(t, l3, []int{4, 5})
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	checkList(t, l1, []int{1, 2, 3})
	checkList(t, l2, []int{4, 5})

	l3.Clear()
	l3.PushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushBackList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3.Clear()
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3.Clear()
	l1.PushBackList(l3)
	checkList(t, l1, []int{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []int{1, 2, 3})
}

func TestRemove(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.InsertAfter(2, e1)
	e3 := l.InsertAfter(3, e2)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3})
	l.RemoveFunc(func(i int) bool { return i&1 != 0 })
	checkListPointers(t, l, []*Element[int]{e2})
	e1 = l.PushFront(1)
	e3 = l.PushBack(3)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3})
	//
	l2 := New[int]()
	e4 := l2.PushFront(4)
	l.InsertAfter(4, e4)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3})
	l.RemoveAfter(e4)
	e4 = l.InsertAfter(4, e3)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})
	l.RemoveFunc(func(i int) bool { return true })
	checkListPointers(t, l, []*Element[int]{})
}

func TestPop(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	checkListPointers(t, l, []*Element[int]{e1, e2})
	//
	l2 := New[int]()
	v, ok := l.PopFront()
	for ok {
		l2.PushBack(v)
		v, ok = l.PopFront()
	}
	checkListPointers(t, l, []*Element[int]{})
	checkList(t, l2, []int{1, 2})
	e := l2.Front()
	v, ok = l2.PopFront()
	if v != e1.Value {
		t.Errorf("e.value = %d, want %v", e.Value, v)
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
}

func TestReverse(t *testing.T) {
	l := New[int]()
	l.Reverse()
	checkList(t, l, []int{})
	l.PushBack(1)
	l.Reverse()
	checkList(t, l, []int{1})
	l.PushBack(2)
	l.Reverse()
	checkList(t, l, []int{2, 1})
	l.Reverse()
	l.PushBack(3)
	checkList(t, l, []int{1, 2, 3})
	l.Reverse()
	checkList(t, l, []int{3, 2, 1})
	l.Reverse()
	l.PushBack(4)
	checkList(t, l, []int{1, 2, 3, 4})
	l.Reverse()
	checkList(t, l, []int{4, 3, 2, 1})
	l.Reverse()
	checkList(t, l, []int{1, 2, 3, 4})

}

func TestUnique(t *testing.T) {
	eq := func(a, b int) bool { return a == b }
	l := New[int]()
	l.Unique(eq)
	checkList(t, l, []int{})
	e1 := l.PushBack(1)
	checkListPointers(t, l, []*Element[int]{e1})
	l.PushBack(1)
	l.Unique(eq)
	checkListPointers(t, l, []*Element[int]{e1})
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	l.PushBack(3)
	l.Unique(eq)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3})
	e4 := l.PushBack(1)
	l.Unique(eq)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})

}

func TestDo(t *testing.T) {
	l := New[int]()
	l.PushFront(1)
	l.PushFront(0)
	l.PushFront(2)
	l.PushFront(3)

	sum := 0
	f := func(v int) {
		sum += v
	}
	l.Do(f)
	if sum != 6 {
		t.Errorf("sum got=%d want 6", sum)
	}
}
