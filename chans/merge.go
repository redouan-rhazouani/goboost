//Package chans implements various algorithms on channels
package chans

func Drain[T any](c <-chan T) {
	for range c {
	}
}

// Merge two channels of some element type into a single element
func Merge[T any](c1, c2 <-chan T) <-chan T {
	r := make(chan T)
	go func(c1, c2 <-chan T, r chan<- T) {
		defer close(r)
		for c1 != nil || c2 != nil {
			select {
			case v, ok := <-c1:
				if ok {
					r <- v
				}
				c1 = nil
			case v, ok := <-c2:
				if ok {
					r <- v
				}
				c2 = nil
			}
		}

	}(c1, c2, r)
	return r
}
