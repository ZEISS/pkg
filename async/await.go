package async

import (
	"sync"
)

// New returns a new future.
func New[T any](fn func() (T, error)) *Future[T] {
	future := Future[T]{value: make(chan T)}

	go func() { // todo: use a sync.Pool for the goroutine
		value, err := fn()
		future.err = err
		future.value <- value

		close(future.value)

		if future.onthen != nil {
			future.onthen(value)
		}

		if future.oncatch != nil {
			future.oncatch(err)
		}
	}()

	return &future
}

// All returns a future that resolves when all the provided futures resolve.
func All[T any](futures ...*Future[T]) *Future[[]T] {
	return New(func() ([]T, error) {
		var (
			w          sync.WaitGroup
			syncValues sync.Mutex
			errOnce    sync.Once
			err        error
		)

		values := make([]T, len(futures))

		for idx, future := range futures {
			w.Add(1)

			go func(future *Future[T], i int) {
				defer w.Done()

				v, err := future.Await()
				if err != nil {
					errOnce.Do(func() {
						future.err = err
					})
					return
				}

				syncValues.Lock()
				defer syncValues.Unlock()

				values[i] = v
			}(future, idx)
		}

		w.Wait()

		return values, err
	})
}
