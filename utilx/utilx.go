package utilx

// Empty returns true of the given value is equal to the default empty value.
func Empty[T comparable](value T) bool {
	var empty T

	return value == empty
}

// NotEmpty returns true of the given value is not equal to the default empty value.
func NotEmpty[T comparable](value T) bool {
	return !Empty(value)
}

// And works similar to "&&" in other languages.
func And[T comparable](a, b T) T {
	var c T

	if a == c {
		return a
	}

	return b
}

// Or works similar to "||" in other languages.
func Or[T comparable](a, b T) T {
	var c T

	if a != c {
		return a
	}

	return b
}

// IfElse works similar to "?:", "if-else" in other languages.
func IfElse[T any](condition bool, v1, v2 T) T {
	if condition {
		return v1
	}

	return v2
}
