package concurrency

import (
	"context"
	"sync"
	"sync/atomic"
)

// TaskResult holds the result of a computation
type TaskResult[T any] struct {
	Result T
	Err    error
}

// Task represents a preemtive function
type Task[T any] func(context.Context) (T, error)

// RunGroup runs a functions concurrently using a sync.WaitGroup
func RunGroup(fns ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		go func(f func()) {
			defer wg.Done()
			f()
		}(fn)
	}
	wg.Wait()
}

// WhenAll runs tasks concurrently and wait until
// all tasks finish successfuly or a leat one of them fails.
// The returned channel contains the result of each task.
func WhenAll[T any](ctx context.Context, tasks ...Task[T]) <-chan TaskResult[T] {
	N := len(tasks)
	results := make(chan TaskResult[T], N)
	go func(results chan<- TaskResult[T]) {
		defer close(results)

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		// run tasks
		errs := make(chan error, N)
		for _, f := range tasks {
			go func(f Task[T]) {
				res, err := f(ctx)
				results <- TaskResult[T]{res, err}
				errs <- err
			}(f)
		}
		// wait for all task to stop
		for i := 0; i < N; i++ {
			select {
			case err := <-errs:
				if err != nil {
					cancel() // it is safe to call cancel many times
				}
			}
		}
	}(results)
	return results
}

// WhenAny runs tasks concurrently and returns when any task executes successfully
func WhenAny[T any](ctx context.Context, funcs ...Task[T]) <-chan TaskResult[T] {
	N := len(funcs)
	results := make(chan TaskResult[T], 1)
	go func(results chan<- TaskResult[T]) {
		defer close(results)
		var count int32
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		errs := make(chan error, N)
		for _, fn := range funcs {
			go func(f Task[T]) {
				res, err := f(ctx)
				if err == nil && atomic.CompareAndSwapInt32(&count, 0, 1) {
					results <- TaskResult[T]{res, err}
					cancel() // winner
				}
				errs <- err
			}(fn)
		}
		// wait for all goroutines to stop
		var firstErr error
		for i := 0; i < N; i++ {
			if err := <-errs; err != nil && firstErr == nil {
				firstErr = err
			}
		}

		// all tasks were not successful
		if atomic.LoadInt32(&count) < 1 {
			results <- TaskResult[T]{Err: firstErr}
		}
	}(results)
	return results
}
