package sync

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// A Mutex is a mutual exclusion lock.
// In additon to standard lock interface Mutex allows goroutines to
// cancell waiting for mutex to be availale when their contexts get cancelled.
// The use need to call NewMutex() to get an unlocked mutex
// A Mutex must not be copied after first use.

type Mutex semaphore.Weighted

// NewMutex creates a new mutex that is ready for use
func NewMutex() *Mutex {
	return (*Mutex)(semaphore.NewWeighted(1))
}

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
	(*semaphore.Weighted)(m).Acquire(context.Background(), 1)
}

// Unlock unlocks m.
// It will panic if m has been locked before
func (m *Mutex) Unlock() {
	(*semaphore.Weighted)(m).Release(1)
}

// TryLock tries to lock m and reports whether it succeeded
//
// Note that while correct uses of TryLock do exist, they are rare,
// and use of TryLock is often a sign of a deeper problem
// in a particular use of mutexes.
func (m *Mutex) TryLock() bool {
	return (*semaphore.Weighted)(m).TryAcquire(1)
}

// LockWithContext tries to lock m, blocking until mutex
// is available or ctx is done. On success, returns nil. On failure, returns
// ctx.Err()
//
// If ctx is already done, LockWithContext may still succeed without blocking.
func (m *Mutex) LockWithContext(ctx context.Context) error {
	return (*semaphore.Weighted)(m).Acquire(ctx, 1)
}
