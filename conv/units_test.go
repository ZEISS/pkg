package conv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/conv"
)

func TestByteSizes(t *testing.T) {
	tests := []struct {
		b    float64
		size string
		want int64
	}{
		{1, "B", 1},
		{1, "KB", 1024},
		{1, "MB", 1048576},
		{1, "GB", 1073741824},
		{1, "TB", 1099511627776},
		{1, "PB", 1125899906842624},
		{1, "EB", 1152921504606846976},
		{1, "ZB", 0},
		{10, "B", 10},
		{10, "KB", 10240},
		{10, "MB", 10485760},
		{10, "GB", 10737418240},
		{10, "TB", 10995116277760},
		{10, "PB", 11258999068426240},
		{1000, "B", 1000},
		{1000, "KB", 1024000},
		{1000, "MB", 1048576000},
		{1000, "GB", 1073741824000},
		{1000, "TB", 1099511627776000},
		{1000, "PB", 1125899906842624000},
		{2.5, "MB", 2621440},
		{2.5, "GB", 2684354560},
		{2.5, "TB", 2748779069440},
		{2.5, "PB", 2814749767106560},
	}

	for _, tt := range tests {
		got := conv.ByteSizes(tt.b, tt.size)
		assert.Equal(t, tt.want, got)
	}
}
