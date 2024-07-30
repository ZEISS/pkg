package errorx

// Ignore is a helper function to ignore errors.
func Ignore[T any](val T, err error) T {
	return val
}
