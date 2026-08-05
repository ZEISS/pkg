package homedir

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	home, err := Get()
	require.NoError(t, err)
	require.NotEmpty(t, home)
	require.True(t, filepath.IsAbs(home))
}
