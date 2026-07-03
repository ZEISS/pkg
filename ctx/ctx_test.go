package ctx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	defaultVal := "default"

	key := New("myKey", defaultVal)

	require.NotNil(t, key)
	require.NotNil(t, key.val)
	require.Equal(t, defaultVal, *key.val)
}

func TestWithValue(t *testing.T) {
	key := New("myKey", "default")
	ctx := key.WithValue(t.Context(), "myValue")

	val, ok := key.ValueOk(ctx)

	require.True(t, ok)
	require.Equal(t, "myValue", val)
}
