package cast

// Slice casts any interface to a slice.
func Slice[T any](val interface{}) []T {
	return val.([]T)
}
