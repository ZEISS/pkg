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
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	require.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	old := strings.Join([]string{tempDir, "example.txt"}, "/")

	f, err := os.Create(old)
	require.NoError(t, err)

	oldBytes, err := f.WriteString("Hello World")
	require.NoError(t, err)
	f.Close()

	new := strings.Join([]string{tempDir, "example_copy.txt"}, "/")

	newBytes, err := filex.CopyFile(old, new, false)
	require.NoError(t, err)

	assert.Equal(t, oldBytes, int(newBytes))

	b, err := os.ReadFile(new)
	require.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}

func TestCopyFileHomeDir(t *testing.T) {
	sr, err := user.Current()
	require.NoError(t, err)

	tempDir, err := os.MkdirTemp(sr.HomeDir, "empty_test")
	require.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	old := strings.Join([]string{tempDir, "example.txt"}, "/")

	f, err := os.Create(old)
	require.NoError(t, err)

	oldBytes, err := f.WriteString("Hello World")
	require.NoError(t, err)
	f.Close()

	new := strings.Join([]string{tempDir, "example_copy.txt"}, "/")
	newHomeDir := strings.Replace(new, sr.HomeDir, "~", 1)

	newBytes, err := filex.CopyFile(old, newHomeDir, false)
	require.NoError(t, err)

	assert.Equal(t, oldBytes, int(newBytes))

	b, err := os.ReadFile(new)
	require.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}

func TestCopyFileMkdir(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "empty_test")
	require.NoError(t, err)

	defer func() { _ = os.RemoveAll(tempDir) }()

	old := strings.Join([]string{tempDir, "example.txt"}, "/")

	f, err := os.Create(old)
	require.NoError(t, err)

	oldBytes, err := f.WriteString("Hello World")
	require.NoError(t, err)
	f.Close()

	new := strings.Join([]string{tempDir, "whoopsy", "example_copy.txt"}, "/")

	newBytes, err := filex.CopyFile(old, new, true)
	require.NoError(t, err)

	assert.Equal(t, oldBytes, int(newBytes))

	b, err := os.ReadFile(new)
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
