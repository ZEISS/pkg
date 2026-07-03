package errorx

// Ignore is a helper function to ignore errors.
func Ignore[T any](val T, err error) T {
	return val
}

// Nil is a helper function to return a nil error.
func Nil[T any](val T) error {
	return nil
}
