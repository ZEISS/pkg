package cast_test

import (
	"testing"

	"github.com/zeiss/pkg/cast"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    any
	}{
		{"success", 1},
		{"success", "hello"},
		{"success", struct{}{}},
		{"success", []int{1, 2, 3}},
		{"success", map[string]int{"a": 1, "b": 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &tt.v
			got := cast.Ptr(tt.v)
			assert.Equal(t, p, got)
		})
	}
}

func TestValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    any
	}{
		{"success", 1},
		{"success", "hello"},
		{"success", struct{}{}},
		{"success", []int{1, 2, 3}},
		{"success", map[string]int{"a": 1, "b": 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := cast.Value(&tt.v)
			assert.Equal(t, tt.v, got)
		})
	}
}
