package cast

// PtrSlice casts an interface to a slice of pointers.
func PtrSlice[T any](slice ...T) []*T {
	ps := make([]*T, 0, len(slice))

	for _, e := range slice {
		ps = append(ps, &e)
	}

	return ps
}
