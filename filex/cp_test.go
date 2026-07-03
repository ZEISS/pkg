package filex_test

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeiss/pkg/filex"
)

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	defer func() { _ = os.RemoveAll(tempDir) }()

	old := strings.Join([]string{tempDir, "example.txt"}, "/")

	f, err := os.Create(old)
	require.NoError(t, err)

	oldBytes, err := f.WriteString("Hello World")
	require.NoError(t, err)
	f.Close()

	newFile := strings.Join([]string{tempDir, "example_copy.txt"}, "/")

	newBytes, err := filex.CopyFile(old, newFile, false)
	require.NoError(t, err)

	assert.Equal(t, oldBytes, int(newBytes))

	b, err := os.ReadFile(newFile)
	require.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}

func TestCopyFileHomeDir(t *testing.T) {
	sr, err := user.Current()
	require.NoError(t, err)

	tempDir := t.TempDir()

	defer func() { _ = os.RemoveAll(tempDir) }()

	old := strings.Join([]string{tempDir, "example.txt"}, "/")

	f, err := os.Create(old)
	require.NoError(t, err)

	oldBytes, err := f.WriteString("Hello World")
	require.NoError(t, err)
	f.Close()

	newFile := strings.Join([]string{tempDir, "example_copy.txt"}, "/")
	newHomeDir := strings.Replace(newFile, sr.HomeDir, "~", 1)

	newBytes, err := filex.CopyFile(old, newHomeDir, false)
	require.NoError(t, err)

	assert.Equal(t, oldBytes, int(newBytes))

	b, err := os.ReadFile(newFile)
	require.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}

func TestCopyFileMkdir(t *testing.T) {
	tempDir := t.TempDir()

	defer func() { _ = os.RemoveAll(tempDir) }()

	old := strings.Join([]string{tempDir, "example.txt"}, "/")

	f, err := os.Create(old)
	require.NoError(t, err)

	oldBytes, err := f.WriteString("Hello World")
	require.NoError(t, err)
	f.Close()

	newStr := strings.Join([]string{tempDir, "whoopsy", "example_copy.txt"}, "/")

	newBytes, err := filex.CopyFile(old, newStr, true)
	require.NoError(t, err)

	assert.Equal(t, oldBytes, int(newBytes))

	b, err := os.ReadFile(newStr)
	require.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}

func TestPrependHomeFolder(t *testing.T) {
	sr, err := user.Current()
	require.NoError(t, err)

	tests := []struct {
		desc        string
		path        string
		expected    string
		expectedErr error
	}{
		{
			path:        "csync/.nanorc",
			expected:    filepath.Join(sr.HomeDir, "csync/.nanorc"),
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			p, err := filex.PrependHomeFolder(tc.path)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, p)
		})
	}
}
