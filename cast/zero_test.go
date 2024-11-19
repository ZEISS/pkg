package cast_test

import (
	"testing"

	"github.com/zeiss/pkg/cast"

	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	assert.Equal(t, 0, cast.Zero[int]())
	assert.Equal(t, "", cast.Zero[string]())
	assert.False(t, cast.Zero[bool]())
	assert.Equal(t, struct{}{}, cast.Zero[struct{}]())
}

func TestIsZero(t *testing.T) {
	assert.True(t, cast.IsZero(0))
	assert.True(t, cast.IsZero(""))
	assert.True(t, cast.IsZero(false))
	assert.True(t, cast.IsZero(struct{}{}))
}
