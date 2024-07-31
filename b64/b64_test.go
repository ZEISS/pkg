package b64_test

import (
	"testing"

	"github.com/zeiss/pkg/b64"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"success", "hello", "aGVsbG8="},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := b64.Base64(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}
