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

// Equal returns true of the given values are equal.
func Equal[T comparable](a, b T) bool {
	return IfElse(Empty(a) && Empty(b), true, a == b)
}

// NotEqual returns true of the given values are not equal.
func NotEqual[T comparable](a, b T) bool {
	return !Equal(a, b)
}

// IsNil returns true of the given value is nil.
func IsNil(value any) bool {
	return value == nil
}

// NotNil returns true of the given value is not nil.
func NotNil(value any) bool {
	return !IsNil(value)
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
