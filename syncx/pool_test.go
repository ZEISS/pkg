package syncx_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/syncx"
)

func TestPoo(t *testing.T) {
	pool := syncx.NewPool(func() *bytes.Buffer {
		return bytes.NewBuffer(make([]byte, 0, 25))
	})

	want := "HELLO"

	var buff *bytes.Buffer = pool.Get()
	assert.Equal(t, 25, buff.Cap())
	assert.Equal(t, 0, buff.Len())

	buff.Reset()
	buff.WriteString(want)
	pool.Put(buff)
}

func BenchmarkPool(b *testing.B) {
	pool := syncx.NewPool(func() *bytes.Buffer {
		return bytes.NewBuffer(make([]byte, 0, 25))
	})

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buff := pool.Get()
			buff.Reset()
			pool.Put(buff)
		}
	})
}
