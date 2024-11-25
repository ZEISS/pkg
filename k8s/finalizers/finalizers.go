package finalizers

import "github.com/zeiss/pkg/slices"

// Finalizer is an interface that can be used to add, remove and check for
type Finalizer interface {
	GetFinalizers() []string
}

// AddFinalizer is a helper function to add a finalizer
func AddFinalizer(obj Finalizer, finalizer string) []string {
	finalizers := obj.GetFinalizers()

	for _, f := range finalizers {
		if finalizer == f {
			return finalizers
		}
	}

	return slices.Append([]string{finalizer}, finalizers...)
}

// HasFinalizer is a helper function to check if a finalizer is present
func HasFinalizer(obj Finalizer, finalizer string) bool {
	finalizers := obj.GetFinalizers()

	for _, f := range finalizers {
		if finalizer == f {
			return true
		}
	}

	return false
}

// RemoveFinalizer is a helper function to remove a finalizer
func RemoveFinalizer(obj Finalizer, finalizer string) []string {
	finalizers := []string{}

	for _, f := range obj.GetFinalizers() {
		if f == finalizer {
			continue
		}
		finalizers = slices.Append(finalizers, f)
	}

	return finalizers
}
