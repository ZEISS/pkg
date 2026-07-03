package filex

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMkdirAll(t *testing.T) {
	tempDir := t.TempDir()
	defer func() { _ = os.RemoveAll(tempDir) }()

	path := strings.Join([]string{tempDir, "example"}, "/")

	err := MkdirAll(path, 0o755)
	require.NoError(t, err)

	_, err = os.Stat(path)
	require.NoError(t, err)
}
