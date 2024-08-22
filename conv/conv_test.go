package conv_test

import (
	"testing"

	"github.com/zeiss/pkg/conv"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  interface{}
	}{
		{
			name: "bool",
			in:   "true",
			out:  true,
		},
		{
			name: "int",
			in:   "1",
			out:  1,
		},
		{
			name: "int8",
			in:   "1",
			out:  int8(1),
		},
		{
			name: "int16",
			in:   "1",
			out:  int16(1),
		},
		{
			name: "int32",
			in:   "1",
			out:  int32(1),
		},
		{
			name: "int64",
			in:   "1",
			out:  int64(1),
		},
		{
			name: "uint",
			in:   "1",
			out:  uint(1),
		},
		{
			name: "uint8",
			in:   "1",
			out:  uint8(1),
		},
		{
			name: "uint16",
			in:   "1",
			out:  uint16(1),
		},
		{
			name: "uint32",
			in:   "1",
			out:  uint32(1),
		},
		{
			name: "uint64",
			in:   "1",
			out:  uint64(1),
		},
		{
			name: "float32",
			in:   "1.000000",
			out:  float32(1),
		},
		{
			name: "float64",
			in:   "1.000000",
			out:  float64(1),
		},
		{
			name: "string",
			in:   "hello",
			out:  "hello",
		},
		{
			name: "nil",
			in:   "",
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := conv.String(tt.out)
			assert.Equal(t, tt.in, got)
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name string
		in   bool
		out  interface{}
	}{
		{
			name: "bool",
			in:   true,
			out:  true,
		},
		{
			name: "int",
			in:   true,
			out:  1,
		},
		{
			name: "int8",
			in:   true,
			out:  int8(1),
		},
		{
			name: "int16",
			in:   true,
			out:  int16(1),
		},
		{
			name: "int32",
			in:   true,
			out:  int32(1),
		},
		{
			name: "int64",
			in:   true,
			out:  int64(1),
		},
		{
			name: "uint",
			in:   true,
			out:  uint(1),
		},
		{
			name: "uint8",
			in:   true,
			out:  uint8(1),
		},
		{
			name: "uint16",
			in:   true,
			out:  uint16(1),
		},
		{
			name: "uint32",
			in:   true,
			out:  uint32(1),
		},
		{
			name: "uint64",
			in:   true,
			out:  uint64(1),
		},
		{
			name: "float32",
			in:   true,
			out:  float32(1),
		},
		{
			name: "float64",
			in:   true,
			out:  float64(1),
		},
		{
			name: "string",
			in:   true,
			out:  "hello",
		},
		{
			name: "nil",
			in:   false,
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := conv.Bool(tt.out)
			assert.Equal(t, tt.in, got)
		})
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		out  interface{}
	}{
		{
			name: "bool",
			in:   []byte("true"),
			out:  true,
		},
		{
			name: "int",
			in:   []byte("1"),
			out:  1,
		},
		{
			name: "int8",
			in:   []byte("1"),
			out:  int8(1),
		},
		{
			name: "int16",
			in:   []byte("1"),
			out:  int16(1),
		},
		{
			name: "int32",
			in:   []byte("1"),
			out:  int32(1),
		},
		{
			name: "int64",
			in:   []byte("1"),
			out:  int64(1),
		},
		{
			name: "uint",
			in:   []byte("1"),
			out:  uint(1),
		},
		{
			name: "uint8",
			in:   []byte("1"),
			out:  uint8(1),
		},
		{
			name: "uint16",
			in:   []byte("1"),
			out:  uint16(1),
		},
		{
			name: "uint32",
			in:   []byte("1"),
			out:  uint32(1),
		},
		{
			name: "uint64",
			in:   []byte("1"),
			out:  uint64(1),
		},
		{
			name: "float32",
			in:   []byte("1.000000"),
			out:  float32(1),
		},
		{
			name: "float64",
			in:   []byte("1.000000"),
			out:  float64(1),
		},
		{
			name: "string",
			in:   []byte("hello"),
			out:  "hello",
		},
		{
			name: "nil",
			in:   []byte(""),
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := conv.Bytes(tt.out)
			assert.Equal(t, tt.in, got)
		})
	}
}
