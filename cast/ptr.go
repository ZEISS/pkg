package cast

// Ptr returns a pointer to the value.
func Ptr[T any](val T) *T {
	return &val
}
