package syncx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/syncx"
)

func TestValue(t *testing.T) {
	ptr := cast.Ptr("hello")

	v := syncx.Value(ptr)
	v.Set(func(value *string) {
		*value = "world"
	})

	assert.Equal(t, "world", cast.Value(ptr))
}

func BenchmarkValue(b *testing.B) {
	ptr := cast.Ptr("hello")

	v := syncx.Value(ptr)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.Set(func(value *string) {
				*value = "world"
			})
		}
	})
}
