package mapx

import (
	"github.com/zeiss/pkg/slices"
)

// Delete removes elements from a map by key.
func Delete[T1 comparable, T2 any](m map[T1]T2, keys ...T1) {
	for _, k := range keys {
		delete(m, k)
	}
}

// Keep removes elements from a map by key.
func Keep[T1 comparable, T2 any](m map[T1]T2, keys ...T1) {
	for k := range m {
		if !slices.In(k, keys...) {
			delete(m, k)
		}
	}
}
