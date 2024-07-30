package cast

// Zero returns a zero value of the given type.
func Zero[T any]() T {
	return *new(T)
}

// IsZero returns true if the given value is the zero value of its type.
func IsZero[T comparable](v T) bool {
	return v == *new(T)
}
