//Package forward_list implements a singly-linked list
//
// To iterate over a list (where l is a *ForwardList):
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
//
package forward_list

type Element[T any] struct {
	// Next pointer to the next element in the singly-linked list.
	next *Element[T]
	// The list to which this element belong.
	list *ForwardList[T]
	// The value stored with this element.
	Value T
}

// Next returns next list element or nil
func (e *Element[T]) Next() *Element[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// ForwardList is a singly-linked list.
// New() is required to create a ready to use empty list.
type ForwardList[T any] struct {
	root Element[T]  // sentinel list element
	back *Element[T] // pointer to the last element in the list
	len  int         // current list length excluding sentinel element
}

// New returns a new initialized singly-linked list
func New[T any]() *ForwardList[T] {
	return new(ForwardList[T]).init()
}

// init initializes or clears list l.
func (l *ForwardList[T]) init() *ForwardList[T] {
	l.root.next = &l.root
	l.root.list = l
	l.back = &l.root
	l.len = 0
	return l
}

// Clear clears all elements from list l.
// After this call Len() returns zero.
func (l *ForwardList[T]) Clear() {
	l.init()
}

// Len returns number of elements of list l
// Complexity O(1)
func (l *ForwardList[T]) Len() int {
	return l.len
}

// Front returns the first element of list l or nil if list is empty
func (l *ForwardList[T]) Front() *Element[T] {
	if l.len > 0 {
		return l.root.next
	}
	return nil
}

// Back returns the last element of list l or nil if list is empty
func (l *ForwardList[T]) Back() *Element[T] {
	if l.len > 0 {
		return l.back

	}
	return nil
}

// insert inserts e after at, increments l.len, and returns e
func (l *ForwardList[T]) insert(e, at *Element[T]) *Element[T] {
	e.next = at.next
	at.next = e
	e.list = l
	if e.next == &l.root {
		l.back = e
	}
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value:v}, at)
func (l *ForwardList[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{Value: v}, at)
}

// PushBack inserts a new element e with value v at the front of the list l and return e
func (l *ForwardList[T]) PushFront(val T) *Element[T] {
	return l.insertValue(val, &l.root)
}

// PushBack inserts a new element e with value v at the back of the list l and return e
func (l *ForwardList[T]) PushBack(v T) *Element[T] {
	return l.insertValue(v, l.back)
}

// InsertAfter inserts a new Element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil
func (l *ForwardList[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark)
}

// Do invokes f on each element of l
// It's not ok to call delete method
func (l ForwardList[T]) Do(f func(T)) {
	for e := l.Front(); e != nil; e = e.Next() {
		f(e.Value)
	}
}

// // Remove element e from l if e is an element of list l.
// // It returns e.next or nil.
// // The element must not be nil
// // complexcity O(n)
// func (l *ForwardList[T]) Remove(e *Element[T]) *Element[T] {
// 	if e.list == l {
// 		return l.remove(e, l.prev(e))
// 	}
// 	return nil
// }

// Remove element after the element e from the list
// It returns the element e.next.next or nil
// The element e must not be nil
// complexcity O(1)
func (l *ForwardList[T]) RemoveAfter(e *Element[T]) *Element[T] {
	if e.list == l {
		if e := l.remove(e.next, e); e != &l.root {
			return e
		}
	}
	return nil
}

// PopFront removes element from the front of the list l
// complexcity O(1)
func (l *ForwardList[T]) PopFront() (T, bool) {
	var value T
	if l.Len() < 1 {
		return value, false
	}
	e := l.root.next
	value = e.Value
	l.remove(e, &l.root)
	return value, true

}

// RemoveFunc removes all elements satisfying f()
func (l *ForwardList[T]) RemoveFunc(f func(v T) bool) {
	prev := &l.root
	for e := l.Front(); e != nil; {
		if f(e.Value) {
			e = l.RemoveAfter(prev)
			continue
		}
		prev = e
		e = e.Next()
	}
}

// remove removes e from list and decrements l.len
// It returns e.next or &l.root
// The element must not be nil
func (l *ForwardList[T]) remove(e, prev *Element[T]) *Element[T] {
	prev.next = e.next
	if e == l.back {
		l.back = prev
	}
	e.next = nil
	l.len--
	//fmt.Printf("%p %p %p\n", &l.root, l.root.next, l.back)
	return prev.next
}

// // prev returns previous element or nil
// func (l *ForwardList[T]) prev(e *Element[T]) *Element[T] {
// 	prev := &l.root
// 	for n := l.Front(); n != nil && n != e; n = n.Next() {
// 		prev = n
// 	}
// 	return prev
// }

// PushBackList inserts a copy of another list at the back of list l.
// The list l and other may by the same.
// They must not be nil
func (l *ForwardList[T]) PushBackList(other *ForwardList[T]) {
	for n, e := other.Len(), other.Front(); n > 0; n, e = n-1, e.Next() {
		l.PushBack(e.Value)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The list l and other may be the same.
// They must not be nil
func (l *ForwardList[T]) PushFrontList(other *ForwardList[T]) {
	next := &l.root
	for n, e := other.Len(), other.Front(); n > 0; n, e = n-1, e.Next() {
		next = l.insertValue(e.Value, next)
	}
}

// Reverse reverses the order of the elements in list l
func (l *ForwardList[T]) Reverse() {
	if l.Len() < 2 {
		return
	}
	first := l.Front()
	last := l.Back()
	reverse(first, nil)
	l.back = first
	l.root.next = last
}

func reverse[T any](first, end *Element[T]) {
	var prev *Element[T]
	for e := first; e != end; {
		e.next, e, prev = prev, e.next, e
	}
}

// Unique removes all consecutive duplicate elements from list l.
// Only the first element in each group of equal elements is left.
func (l *ForwardList[T]) Unique(f func(a, b T) bool) {
	if l.Len() < 1 {
		return
	}
	prev := l.Front()
	for e := prev.Next(); e != nil; {
		if f(prev.Value, e.Value) {
			e = l.RemoveAfter(prev)
			continue
		}
		prev = e
		e = e.Next()
	}
}
