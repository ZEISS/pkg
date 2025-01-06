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

func TestExists(t *testing.T) {
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	assert.True(t, mapx.Exists(m, 2))
	assert.False(t, mapx.Exists(m, 4))
}

func TestMerge(t *testing.T) {
	m1 := map[int]string{
		1: "one",
		2: "two",
	}
	m2 := map[int]string{
		3: "three",
		4: "four",
	}
	m3 := map[int]string{
		4: "foo",
		5: "five",
		6: "six",
	}
	assert.Equal(t, map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "foo",
		5: "five",
		6: "six",
	}, mapx.Merge(m1, m2, m3))
}
