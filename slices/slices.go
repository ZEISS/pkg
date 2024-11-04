package slices

import (
	"github.com/zeiss/pkg/cast"
)

// Any checks if any element in a slice satisfies a predicate.
func Any[T any](fn func(v T) bool, slice ...T) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}

	return false
}

// Limit returns a slice with a maximum length of limit.
func Limit[T any](limit int, slice ...T) []T {
	if len(slice) > limit {
		return slice[:limit]
	}

	return slice
}

// Map applies a function to all elements in a slice.
func Map[T1 any, T2 any](fn func(v T1) T2, slice ...T1) []T2 {
	a := make([]T2, len(slice))
	for i, v := range slice {
		a[i] = fn(v)
	}

	return a
}

// Range returns a slice with elements from start to end.
func Range(from, to int) []int {
	result := make([]int, to-from+1)
	for i := 0; i <= to-from; i++ {
		result[i] = i + from
	}

	return result
}

// Slice casts an interface to a slice of a specific type.
func Slice[T any](slice []any) []T {
	newslice := make([]T, 0, len(slice))

	for _, el := range slice {
		newslice = append(newslice, el.(T)) //nolint:forcetypeassert // We expect a panic, if something is wrong
	}

	return newslice
}

// Cut removes an element from a slice at a given position.
func Cut[T comparable](i, j int, a ...T) []T {
	copy(a[i:], a[j:])
	for k, n := len(a)-j+i, len(a); k < n; k++ {
		a[k] = cast.Zero[T]()
	}

	return a[:len(a)-j+i]
}

// Delete removes an element from a slice by value.
func Delete[T comparable](i int, a ...T) []T {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = cast.Zero[T]()

	return a[:len(a)-1]
}

// Push adds an element to the end of a slice.
func Push[T comparable](x T, a ...T) []T {
	return append(a, x)
}

// Pop removes an element from the end of a slice.
func Pop[T comparable](a ...T) (T, []T) {
	return a[len(a)-1], a[:len(a)-1]
}

// Insert adds an element at a given position in a slice.
func Insert[T comparable](x T, idx int, a ...T) []T {
	return append(a[:idx], append([]T{x}, a[idx:]...)...)
}

// Filter removes all elements from a slice that satisfy a predicate.
func Filter[T any](fn func(v T) bool, slice ...T) []T {
	newslice := make([]T, 0)

	for _, v := range slice {
		if fn(v) {
			newslice = append(newslice, v)
		}
	}

	return newslice
}

// Last returns the last element of a slice.
func Last[T any](slice ...T) T {
	return slice[len(slice)-1]
}

// First returns the first element of a slice.
func First[T any](slice ...T) T {
	return slice[0]
}

// In checks if a value is in a slice.
func In[T comparable](val T, slice ...T) bool {
	m := make(map[T]bool, len(slice))
	for _, v := range slice {
		m[v] = true
	}

	_, ok := m[val]
	return ok
}

// Index returns the index of the first element in a slice that satisfies a predicate.
func Index[T any](fn func(v T) bool, slice ...T) int {
	for i, v := range slice {
		if fn(v) {
			return i
		}
	}

	return -1
}

// Unique returns a slice with all duplicate elements removed.
func Unique[T1 any, T2 comparable](fn func(v T1) T2, slice ...T1) []T1 {
	flags := map[T2]bool{}

	return Filter(func(v T1) bool {
		compareval := fn(v)
		defer func() {
			flags[compareval] = true
		}()

		return !flags[compareval]
	}, slice...)
}

// Size checks if a slice has a specific size.
func Size[T any](size int, slice ...T) bool {
	return size == len(slice)
}
