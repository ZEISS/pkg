package stringx

import "strings"

// FirstN returns the first n characters of a string.
func FirstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}

	return s
}

// AnyPrefix checks if a string has any of the given prefixes.
func AnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if !strings.HasPrefix(s, p) {
			continue
		}

		return true
	}

	return false
}
