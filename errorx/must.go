package errorx

// Must panics if the error is not nil.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

// Panic panics with the given error.
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
