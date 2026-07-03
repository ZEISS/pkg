package filex

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()

	defer func() { _ = os.RemoveAll(tempDir) }()

	ok, err := FileExists(tempDir)
	require.NoError(t, err)
	assert.True(t, ok)

	path := strings.Join([]string{tempDir, "example.txt"}, "/")
	f, err := os.Create(path)
	require.NoError(t, err)
	f.Close()

	ok, err = FileExists(path)
	require.NoError(t, err)
	assert.True(t, ok)

	err = os.Remove(path)
	require.NoError(t, err)

	ok, err = FileExists(path)
	require.Error(t, err)
	assert.False(t, ok)
}

func TestFileNotExists(t *testing.T) {
	tempDir := t.TempDir()
	defer func() { _ = os.RemoveAll(tempDir) }()

	path := strings.Join([]string{tempDir, "example.txt"}, "/")
	f, err := os.Create(path)
	require.NoError(t, err)
	f.Close()

	ok, err := FileNotExists(path)
	require.NoError(t, err)
	assert.False(t, ok)

	err = os.Remove(path)
	require.NoError(t, err)

	ok, err = FileNotExists(path)
	require.NoError(t, err)
	assert.True(t, ok)

	path = strings.Join([]string{tempDir, "demo123"}, "/")
	ok, err = FileNotExists(path)
	require.NoError(t, err)
	assert.True(t, ok)
}
