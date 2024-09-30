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

func TestOr(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 1},
		{0, 2, 2},
	}

	for _, test := range tests {
		got := utilx.Or(test.a, test.b)
		require.Equal(t, test.expected, got)
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 2},
		{0, 2, 0},
	}

	for _, test := range tests {
		got := utilx.And(test.a, test.b)
		require.Equal(t, test.expected, got)
	}
}

func TestNotEmpty(t *testing.T) {
	tests := []struct {
		value    int
		expected bool
	}{
		{0, false},
		{1, true},
	}

	for _, test := range tests {
		got := utilx.NotEmpty(test.value)
		require.Equal(t, test.expected, got)
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		value    int
		expected bool
	}{
		{0, true},
		{1, false},
	}

	for _, test := range tests {
		got := utilx.Empty(test.value)
		require.Equal(t, test.expected, got)
	}
}
