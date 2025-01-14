package reflectx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/reflectx"
)

func TestParseTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tag             string
		expected        string
		expectedOptions reflectx.TagOptions
	}{
		{
			tag:             "name,omitempty",
			expected:        "name",
			expectedOptions: "omitempty",
		},
		{
			tag:             "name,omitempty,required",
			expected:        "name",
			expectedOptions: "omitempty,required",
		},
		{
			tag:             "name,omitempty,required,truncate",
			expected:        "name",
			expectedOptions: "omitempty,required,truncate",
		},
	}

	for _, test := range tests {
		tag, opt := reflectx.ParseTag(test.tag)
		assert.Equal(t, test.expected, tag)
		assert.Equal(t, test.expectedOptions, opt)
	}
}

func TestIsValidTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tag      string
		expected bool
	}{
		{
			tag:      "name",
			expected: true,
		},
		{
			tag:      "",
			expected: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, reflectx.IsValidTag(test.tag))
	}
}
