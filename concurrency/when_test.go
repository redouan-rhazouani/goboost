package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRunGroup(t *testing.T) {
	var a, b int
	RunGroup()
	RunGroup(func() { a = 1 })
	if a != 1 {
		t.Errorf("a got=%d want 1", a)
	}
	RunGroup(func() { a = 1 }, func() { b = 1 })
	if sum := a + b; sum != 2 {
		t.Errorf("a got=%d want 2", a)
	}
}
func TestWhenAllZero(t *testing.T) {
	ctx := context.Background()
	res := WhenAll[int](ctx)
	select {
	case r := <-res:
		if r.Result != 0 || r.Err != nil {
			t.Errorf("Result got=%v want={0, nil}", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}
}

func TestWhenAnyZero(t *testing.T) {
	ctx := context.Background()
	res := WhenAny[int](ctx)
	select {
	case r := <-res:
		if r.Result != 0 || r.Err != nil {
			t.Errorf("Result got=%v want={0, nil}", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}
}

func TestWhenAllOne(t *testing.T) {
	anyError := errors.New("foo")
	ctx := context.Background()
	task := newTask(ctx, "Task-1", 1, time.Millisecond, nil)
	c := WhenAll(ctx, task)
	select {
	case r := <-c:
		if r.Result != 1 || r.Err != nil {
			t.Errorf("Result got=%v want={0, nil}", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}

	task = newTask(ctx, "Task-2", 1, time.Millisecond, anyError)
	c = WhenAll(ctx, task)
	select {
	case r := <-c:
		if r.Result != 0 || !errors.Is(r.Err, anyError) {
			t.Errorf("Result got=%v want={0, %v}", r, anyError)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	task = newTask(ctx, "Task-2", 1, time.Second, nil)
	c = WhenAll(ctx, task)
	select {
	case r := <-c:
		if r.Result != 0 || !errors.Is(r.Err, context.DeadlineExceeded) {
			t.Errorf("Result got=%v want={0, %v}", r, anyError)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}

}

func TestWhenAnyOne(t *testing.T) {
	anyError := errors.New("foo")
	ctx := context.Background()
	task := newTask(ctx, "Task-1", 1, time.Millisecond, nil)
	c := WhenAny(ctx, task)
	select {
	case r := <-c:
		if r.Result != 1 || r.Err != nil {
			t.Errorf("Result got=%v want={0, nil}", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}

	task = newTask(ctx, "Task-2", 1, time.Millisecond, anyError)
	c = WhenAny(ctx, task)
	select {
	case r := <-c:
		if r.Result != 0 || !errors.Is(r.Err, anyError) {
			t.Errorf("Result got=%v want={0, %v}", r, anyError)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	task = newTask(ctx, "Task-2", 1, time.Second, nil)
	c = WhenAny(ctx, task)
	select {
	case r := <-c:
		if r.Result != 0 || !errors.Is(r.Err, context.DeadlineExceeded) {
			t.Errorf("Result got=%v want={0, %v}", r, anyError)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("timeout")
	}

}

func TestWhenAllMany(t *testing.T) {
	errFoo := errors.New("foo")
	ctx := context.Background()
	{
		c := WhenAll(ctx,
			newTask(ctx, "Task-1", 1, time.Millisecond, nil),
			newTask(ctx, "Task-2", 2, time.Millisecond, nil),
			newTask(ctx, "Task-3", 3, time.Millisecond, nil),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil {
				t.Error(err)
			}
		}
		if sum != 6 {
			t.Errorf("sum got=%v want 6", sum)
		}
	}
	{ // one task returns an error
		c := WhenAll(ctx,
			newTask(ctx, "Task-1", 1, time.Millisecond, nil),
			newTask(ctx, "Task-2", 2, time.Millisecond, nil),
			newTask(ctx, "Task-3", 3, time.Millisecond, errFoo),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil && !errors.Is(err, errFoo) &&
				!errors.Is(err, context.Canceled) {
				t.Error(err)
			}
		}
		if sum > 5 {
			t.Errorf("sum got=%v want < 6", sum)
		}
	}
	{ // cancellation
		ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()
		c := WhenAll(ctx,
			newTask(ctx, "Task-1", 1, time.Millisecond, nil),
			newTask(ctx, "Task-2", 2, time.Millisecond, nil),
			newTask(ctx, "Task-3", 3, 100*time.Millisecond, nil),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil && !errors.Is(err, context.DeadlineExceeded) {
				t.Error(err)
			}
		}
		if sum != 3 {
			t.Errorf("sum got=%v want=3", sum)
		}
	}
	{ // all tasks return error
		c := WhenAll(ctx,
			newTask(ctx, "Task-1", 1, time.Millisecond, errFoo),
			newTask(ctx, "Task-2", 2, time.Millisecond, errFoo),
			newTask(ctx, "Task-3", 3, time.Millisecond, errFoo),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil && !errors.Is(err, errFoo) &&
				!errors.Is(err, context.Canceled) {
				t.Error(err)
			}
		}
		if sum != 0 {
			t.Errorf("sum got=%v want=0", sum)
		}
	}
}

func TestWhenAnyMany(t *testing.T) {
	errFoo := errors.New("foo")
	ctx := context.Background()
	{
		c := WhenAny(ctx,
			newTask(ctx, "Task-1", 1, 100*time.Millisecond, nil),
			newTask(ctx, "Task-2", 2, time.Millisecond, nil),
			newTask(ctx, "Task-3", 3, 100*time.Millisecond, nil),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil {
				t.Error(err)
			}
		}
		if sum != 2 {
			t.Errorf("sum got=%v want 6", sum)
		}
	}
	{ // one task returns an error
		c := WhenAny(ctx,
			newTask(ctx, "Task-1", 1, time.Millisecond, errFoo),
			newTask(ctx, "Task-2", 2, time.Millisecond, errFoo),
			newTask(ctx, "Task-3", 3, 3*time.Millisecond, nil),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil {
				t.Error(err)
			}
		}
		if sum != 3 {
			t.Errorf("sum got=%v want=3", sum)
		}
	}
	{ // cancellation
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()
		c := WhenAny(ctx,
			newTask(ctx, "Task-1", 1, 100*time.Millisecond, nil),
			newTask(ctx, "Task-2", 2, 100*time.Millisecond, nil),
			newTask(ctx, "Task-3", 3, 100*time.Millisecond, nil),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil && !errors.Is(err, context.DeadlineExceeded) {
				t.Error(err)
			}
		}
		if sum != 0 {
			t.Errorf("sum got=%v want=0", sum)
		}
	}
	{ // all tasks return error
		c := WhenAny(ctx,
			newTask(ctx, "Task-1", 1, time.Millisecond, errFoo),
			newTask(ctx, "Task-2", 2, time.Millisecond, errFoo),
			newTask(ctx, "Task-3", 3, time.Millisecond, errFoo),
		)
		sum := 0
		for r := range c {
			sum += r.Result
			if err := r.Err; err != nil && !errors.Is(err, errFoo) {
				t.Error(err)
			}
		}
		if sum != 0 {
			t.Errorf("sum got=%v want=0", sum)
		}
	}
}

func newTask(ctx context.Context, name string, ret int, dur time.Duration, err error) Task[int] {
	return func(ctx context.Context) (int, error) {
		select {
		case <-time.After(dur):
		case <-ctx.Done():
			err = ctx.Err()
		}
		if err != nil {
			err = fmt.Errorf("%v err: %w", name, err)
			ret = 0
		}
		return ret, err
	}
}
