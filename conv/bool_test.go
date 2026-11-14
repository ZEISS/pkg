package conv_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/conv"
)

func TestBoolZeroValuesAndPointers(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"int zero", int(0), false},
		{"int nil pointer", (*int)(nil), false},
		{"int zero pointer", cast.Ptr(int(0)), false},
		{"int nonzero pointer", cast.Ptr(int(1)), true},
		{"int8 zero", int8(0), false},
		{"int8 nil pointer", (*int8)(nil), false},
		{"int8 zero pointer", cast.Ptr(int8(0)), false},
		{"int8 nonzero pointer", cast.Ptr(int8(1)), true},
		{"int16 zero", int16(0), false},
		{"int16 nil pointer", (*int16)(nil), false},
		{"int16 zero pointer", cast.Ptr(int16(0)), false},
		{"int16 nonzero pointer", cast.Ptr(int16(1)), true},
		{"int32 zero", int32(0), false},
		{"int32 nil pointer", (*int32)(nil), false},
		{"int32 zero pointer", cast.Ptr(int32(0)), false},
		{"int32 nonzero pointer", cast.Ptr(int32(1)), true},
		{"int64 zero", int64(0), false},
		{"int64 nil pointer", (*int64)(nil), false},
		{"int64 zero pointer", cast.Ptr(int64(0)), false},
		{"int64 nonzero pointer", cast.Ptr(int64(1)), true},
		{"uint zero", uint(0), false},
		{"uint nil pointer", (*uint)(nil), false},
		{"uint zero pointer", cast.Ptr(uint(0)), false},
		{"uint nonzero pointer", cast.Ptr(uint(1)), true},
		{"uint8 zero", uint8(0), false},
		{"uint8 nil pointer", (*uint8)(nil), false},
		{"uint8 zero pointer", cast.Ptr(uint8(0)), false},
		{"uint8 nonzero pointer", cast.Ptr(uint8(1)), true},
		{"uint16 zero", uint16(0), false},
		{"uint16 nil pointer", (*uint16)(nil), false},
		{"uint16 zero pointer", cast.Ptr(uint16(0)), false},
		{"uint16 nonzero pointer", cast.Ptr(uint16(1)), true},
		{"uint32 zero", uint32(0), false},
		{"uint32 nil pointer", (*uint32)(nil), false},
		{"uint32 zero pointer", cast.Ptr(uint32(0)), false},
		{"uint32 nonzero pointer", cast.Ptr(uint32(1)), true},
		{"uint64 zero", uint64(0), false},
		{"uint64 nil pointer", (*uint64)(nil), false},
		{"uint64 zero pointer", cast.Ptr(uint64(0)), false},
		{"uint64 nonzero pointer", cast.Ptr(uint64(1)), true},
		{"float32 zero", float32(0), false},
		{"float32 nil pointer", (*float32)(nil), false},
		{"float32 zero pointer", cast.Ptr(float32(0)), false},
		{"float32 nonzero pointer", cast.Ptr(float32(1)), true},
		{"float64 zero", float64(0), false},
		{"float64 nil pointer", (*float64)(nil), false},
		{"float64 zero pointer", cast.Ptr(float64(0)), false},
		{"float64 nonzero pointer", cast.Ptr(float64(1)), true},
		{"negative integer", int64(-1), true},
		{"negative float", -0.5, true},
		{"negative zero", math.Copysign(0, -1), false},
		{"nan", math.NaN(), true},
		{"infinity", math.Inf(1), true},
		{"false", false, false},
		{"false pointer", cast.Ptr(false), false},
		{"true pointer", cast.Ptr(true), true},
		{"nil bool pointer", (*bool)(nil), false},
		{"empty string", "", false},
		{"empty string pointer", cast.Ptr(""), false},
		{"nonempty string pointer", cast.Ptr("false"), true},
		{"nil string pointer", (*string)(nil), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			if assert.NotPanics(t, func() { got = conv.Bool(tt.value) }) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestBoolsZeroValuesAndPointers(t *testing.T) {
	assert.Equal(t, []bool{false, true, false}, conv.Bools(int64(0), int64(1), int64(0)))
	var value *float64
	assert.NotPanics(t, func() {
		assert.Equal(t, []bool{false, true, false}, conv.Bools(cast.Ptr(0.0), cast.Ptr(0.5), value))
	})
}

func TestBoolUnsupportedType(t *testing.T) {
	assert.Panics(t, func() { conv.Bool(struct{}{}) })
}
