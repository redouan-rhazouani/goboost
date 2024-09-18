package sync

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// Mutex is a mutual exclusion lock with additional support for context cancellation.
// Unlike a standard lock, this Mutex allows goroutines to cancel their wait for the lock
// if their contexts are cancelled, improving responsiveness in cancellation scenarios.
// To create a new, unlocked Mutex, use NewMutex().
// Note: A Mutex must not be copied after it has been used.

type Mutex semaphore.Weighted

// NewMutex creates a new mutex that is ready for use
func NewMutex() *Mutex {
	return (*Mutex)(semaphore.NewWeighted(1))
}

// Lock acquires the mutex 'm'.
// If the mutex is already locked, the calling goroutine will block
// until the lock becomes available
func (m *Mutex) Lock() {
	(*semaphore.Weighted)(m).Acquire(context.Background(), 1)
}

// Unlock releases the mutex 'm'.
// It will panic if 'm' is not locked
func (m *Mutex) Unlock() {
	(*semaphore.Weighted)(m).Release(1)
}

// TryLock attempts to acquire the mutex 'm' without blocking, returning true if
// the lock was successfully acquired, and false otherwise.
//
// While there are legitimate use cases for TryLock, they are uncommon.
// Frequent use of TryLock may indicate a deeper issue in the design or
// usage pattern of mutexes.
func (m *Mutex) TryLock() bool {
	return (*semaphore.Weighted)(m).TryAcquire(1)
}

// LockWithContext attempts to acquire the mutex 'm', blocking until the mutex
// becomes available or the context 'ctx' is done. On success, it returns nil.
// If the operation fails due to 'ctx' being done, it returns ctx.Err().
//
// Note: If 'ctx' is already done, LockWithContext may still acquire the lock without blocking.
func (m *Mutex) LockWithContext(ctx context.Context) error {
	return (*semaphore.Weighted)(m).Acquire(ctx, 1)
}
