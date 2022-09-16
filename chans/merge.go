// Package chans implements various algorithms on channels
package chans

// Drain drains a channel
func Drain[T any](c <-chan T) {
	for range c {
	}
}

// Merge a pair of channels into a single channel
func Merge[T any](c1, c2 <-chan T) <-chan T {
	r := make(chan T)
	go func(c1, c2 <-chan T, r chan<- T) {
		defer close(r)
		for c1 != nil || c2 != nil {
			select {
			case v, ok := <-c1:
				if ok {
					r <- v
				} else {
					c1 = nil
				}
			case v, ok := <-c2:
				if ok {
					r <- v
				} else {
					c2 = nil
				}
			}
		}

	}(c1, c2, r)
	return r
}

// Repeat value v n times
func Repeat[T any](v T, n int) <-chan T {
	r := make(chan T, n)
	go func(c <-chan T, v T, n int) {
		defer close(r)
		for ; n > 0; n-- {
			r <- v
		}
	}(r, v, n)
	return r
}
