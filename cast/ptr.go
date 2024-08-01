package cast

// Ptr returns a pointer to the value.
func Ptr[T any](val T) *T {
	return &val
}

// Value returns the value of the pointer.
func Value[T any](val *T) T {
	if val == nil {
		return Zero[T]()
	}

	return *val
}
