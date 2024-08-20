package urlx_test

import (
	"net/url"
	"testing"

	"github.com/zeiss/pkg/urlx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		values   url.Values
		expected string
	}{
		{
			name:     "copy values",
			input:    "http://example.com",
			values:   url.Values{"key": []string{"value"}},
			expected: "http://example.com?key=value",
		},
		{
			name:     "copy additional values",
			input:    "http://example.com?key=value",
			values:   url.Values{"key2": []string{"value2"}},
			expected: "http://example.com?key=value&key2=value2",
		},
		{
			name:     "overwrite values",
			input:    "http://example.com?key=value",
			values:   url.Values{"key": []string{"value2"}},
			expected: "http://example.com?key=value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := urlx.CopyValues(tt.input, tt.values)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestMustCopyValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		values   url.Values
		expected string
	}{
		{
			name:     "copy values",
			input:    "http://example.com",
			values:   url.Values{"key": []string{"value"}},
			expected: "http://example.com?key=value",
		},
		{
			name:     "copy additional values",
			input:    "http://example.com?key=value",
			values:   url.Values{"key2": []string{"value2"}},
			expected: "http://example.com?key=value&key2=value2",
		},
		{
			name:     "overwrite values",
			input:    "http://example.com?key=value",
			values:   url.Values{"key": []string{"value2"}},
			expected: "http://example.com?key=value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := urlx.MustCopyValues(tt.input, tt.values)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
