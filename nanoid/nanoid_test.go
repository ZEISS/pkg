package nanoid_test

import (
	"testing"

	"github.com/zeiss/pkg/nanoid"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := nanoid.New(1)
	require.Error(t, err)

	_, err = nanoid.New(256)
	require.Error(t, err)
}

func TestUnicode(t *testing.T) {
	_, err := nanoid.Unicode(`🦄🌈🍕💩👽💀🎃🦁🍻🚀`, 10)
	require.NoError(t, err)
}

func Benchmark8NanoID(b *testing.B) {
	f, err := nanoid.New(8)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		f()
	}
}

func Benchmark16NanoID(b *testing.B) {
	f, err := nanoid.New(16)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		f()
	}
}

func Benchmark21NanoID(b *testing.B) {
	f, err := nanoid.New(nanoid.DefaultLength)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		f()
	}
}

func Benchmark32NanoID(b *testing.B) {
	f, err := nanoid.New(32)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		f()
	}
}
