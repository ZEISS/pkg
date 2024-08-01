package cast

// String is a type that represents a string.
func String[T ~string](val T) string {
	return string(val)
}
