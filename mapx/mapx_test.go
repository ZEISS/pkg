package mapx_test

import (
	"testing"

	"github.com/zeiss/pkg/mapx"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	mapx.Delete(m, 2)
	assert.Equal(t, map[int]string{
		1: "one",
		3: "three",
	}, m)
}

func TestKeep(t *testing.T) {
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	mapx.Keep(m, 2)
	assert.Equal(t, map[int]string{
		2: "two",
	}, m)
}
