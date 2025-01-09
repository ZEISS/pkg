package httpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContentTypeByExtension(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{
			name:     "html",
			file:     "index.html",
			expected: "text/html",
		},
		{
			name:     "css",
			file:     "style.css",
			expected: "text/css",
		},
		{
			name:     "js",
			file:     "script.js",
			expected: "application/javascript",
		},
		{
			name:     "unknown",
			file:     "file",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, GetContentTypeByExtension(test.file))
		})
	}
}
