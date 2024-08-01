package utilx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/pkg/utilx"
)

func TestIfElse(t *testing.T) {
	tests := []struct {
		condition bool
		v1        any
		v2        any
		expected  any
	}{
		{true, 1, 2, 1},
		{false, 1, 2, 2},
	}

	for _, test := range tests {
		got := utilx.IfElse(test.condition, test.v1, test.v2)
		require.Equal(t, test.expected, got)
	}
}
