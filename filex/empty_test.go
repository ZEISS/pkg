package filex

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDirEmpty(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	require.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	isEmpty, err := IsDirEmpty(tempDir)
	require.NoError(t, err)
	assert.True(t, isEmpty)

	f, err := os.Create(path.Join(tempDir, "test.txt"))
	require.NoError(t, err)

	f.Close()

	isEmpty, err = IsDirEmpty(tempDir)
	require.NoError(t, err)
	assert.False(t, isEmpty)
}
