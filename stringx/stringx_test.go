package stringx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/stringx"
)

func TestFirstN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"", 0, ""},
		{"", 1, ""},
		{"a", 0, ""},
		{"a", 1, "a"},
		{"a", 2, "a"},
		{"ab", 0, ""},
		{"ab", 1, "a"},
		{"ab", 2, "ab"},
		{"ab", 3, "ab"},
		{"abc", 0, ""},
		{"abc", 1, "a"},
		{"abc", 2, "ab"},
		{"abc", 3, "abc"},
		{"abc", 4, "abc"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%d", tt.s, tt.n), func(t *testing.T) {
			got := stringx.FirstN(tt.s, tt.n)
			assert.Equal(t, tt.want, got)
		})
	}
}
