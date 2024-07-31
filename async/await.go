package async

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
