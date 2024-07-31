package async

import (
	"encoding/json"
	"sync"
)

// Future is a type that represents a future value.
type Future[T any] struct {
	value chan T
	err   error

	onthen  func(T)
	oncatch func(error)

	cache   *T
	runOnce sync.Once
}

// Await  waits for the future to resolve and returns the value.
func Await[T any](f *Future[T]) (T, error) {
	return f.Await()
}

// Await waits for the future to resolve and returns the value.
func (f *Future[T]) Await() (T, error) {
	f.runOnce.Do(func() {
		v := <-f.value
		f.cache = &v
	})

	return *f.cache, f.err
}

// Then sets a callback to be called when the future resolves.
func (f *Future[T]) Then(fn func(T)) *Future[T] {
	f.Await() //nolint:errcheck

	if f.err == nil {
		fn(*f.cache)
	}

	return f
}

// Catch sets a callback to be called when the future rejects.
func (f *Future[T]) Catch(fn func(error)) *Future[T] {
	f.Await() //nolint:errcheck

	if f.err != nil {
		fn(f.err)
	}

	return f
}

// MarshalJSON implements the json.Marshaler interface.
func (f *Future[T]) MarshalJSON() ([]byte, error) {
	val, err := Await(f)
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(val)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Future[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.cache)
}
