package unitx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/pkg/unitx"
)

func TestHumanSize_ToInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		h       unitx.HumanSize
		want    int64
		wantErr bool
	}{
		{
			name:    "zero",
			h:       unitx.HumanSize("0"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1kB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "binary",
			h:       unitx.HumanSize("1KiB"),
			want:    1024,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1M"),
			want:    1000 * 1000,
			wantErr: false,
		},
		{
			name:    "decimal",
			h:       unitx.HumanSize("1MB"),
			want:    1024 * 1024,
			wantErr: false,
		},
		{
			name:    "invalid",
			h:       unitx.HumanSize("1MBa"),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.ToInt()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
