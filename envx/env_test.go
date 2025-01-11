package envx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/envx"
)

func TestEnv(t *testing.T) {
	t.Parallel()

	// Create a new environment.
	e := envx.Env{
		"key1": "value1",
		"key2": "value2",
	}

	e2 := e.Copy()
	assert.Equal(t, e, e2)

	strings := e.Strings()
	assert.Len(t, strings, 2)
	assert.Contains(t, strings, "key1=value1")
	assert.Contains(t, strings, "key2=value2")
}
