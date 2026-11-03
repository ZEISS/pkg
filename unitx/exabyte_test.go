package unitx_test

import (
	"testing"

	"github.com/zeiss/pkg/unitx"
)

func TestHumanSizeExabyteMultipliers(t *testing.T) {
	for _, tc := range []struct {
		input unitx.HumanSize
		want  int64
	}{
		{"1e", 1000000000000000000},
		{"2E", 2000000000000000000},
		{"-1e", -1000000000000000000},
		{"1eb", 1152921504606846976},
		{"1Ei", 1152921504606846976},
		{"1EiB", 1152921504606846976},
		{"2 eIb", 2305843009213693952},
		{"-1eib", -1152921504606846976},
		{"0e", 0},
		{"1p", 1000000000000000},
		{"1PiB", 1125899906842624},
	} {
		t.Run(string(tc.input), func(t *testing.T) {
			got, err := tc.input.ToInt64()
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %d; want %d", got, tc.want)
			}
		})
	}
}
