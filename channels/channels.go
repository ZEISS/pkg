package channels

import "sync"

// Join joins multiple channels into one.
func Join[T any](inputs ...<-chan T) <-chan T {
	c := make(chan T)
	var once sync.Once

	for _, input := range inputs {
		go func(input <-chan T) {
			defer once.Do(func() { close(c) })
			for v := range input {
				c <- v
			}
		}(input)
	}

	return c
}

// Merge merges multiple channels into one.
func Merge[T any](inputs ...<-chan T) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)
		for _, input := range inputs {
			go func(input <-chan T) {
				for v := range input {
					c <- v
				}
			}(input)
		}
	}()

	return c
}

// Broadcast broadcasts a channel to multiple channels.
func Broadcast[T any](input <-chan T, outputs ...chan<- T) {
	go func() {
		for v := range input {
			for _, output := range outputs {
				output <- v
			}
		}
	}()
}

// Drain drains the channel until it is closed.
func Drain[T any](input <-chan T) {
	if input == nil {
		return
	}

	for {
		_, ok := <-input
		if !ok {
			break
		}
	}
}

// Filter filters the channel with the given function.
func Filter[T any](input <-chan T, fn func(T) bool) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)
		for v := range input {
			if fn(v) {
				c <- v
			}
		}
	}()

	return c
}

// Slice slices the channel into a slice.
func Slice[T any](ch <-chan any) []T {
	result := make([]T, 0)
	for e := range ch {
		result = append(result, e.(T))
	}

	return result
}

// Channel is a helper function to send a slice to a channel.
func Channel[T any](source []T, in chan any) {
	for _, e := range source {
		in <- e
	}
}

// Wrap wraps a channel into a typed channel.
func Wrap[T any](ch chan any) chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		for e := range ch {
			out <- e.(T)
		}
	}()

	return out
}
