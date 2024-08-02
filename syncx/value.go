package syncx

import "sync"

// AtomicValue is a strongly typed version of sync/atomic.Value
type AtomicValue[T any] struct {
	mu sync.Mutex
	v  *T
}

// Value returns a new Value
func Value[T any](v *T) *AtomicValue[T] {
	return &AtomicValue[T]{v: v}
}

// Set sets the value
func (v *AtomicValue[T]) Set(fn func(value *T)) {
	v.mu.Lock()
	defer v.mu.Unlock()
	fn(v.v)
}
