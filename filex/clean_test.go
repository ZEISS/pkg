package filex

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	require.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	path := filepath.Join(tempDir, "example")
	err = MkdirAll(path, 0o755)
	require.NoError(t, err)

	file, err := os.Create(filepath.Join(tempDir, "example", "test.txt"))
	require.NoError(t, err)
	defer file.Close()

	err = Clean(path, 0o755)
	require.NoError(t, err)

	_, err = os.Stat(path)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tempDir, "example", "test.txt"))
	require.Error(t, err)
}
