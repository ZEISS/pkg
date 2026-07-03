package envx_test

import (
	"os"
	"path"
	"testing"

	"github.com/zeiss/pkg/envx"

	"github.com/stretchr/testify/require"
)

func TestHasUser(t *testing.T) {
	t.Parallel()

	require.NoError(t, envx.HasUser()(t.Context()))
}

func TestIsDirEmpty(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	defer func() { _ = os.RemoveAll(tempDir) }()

	require.NoError(t, envx.IsDirEmpty(tempDir)(t.Context()))

	f, err := os.Create(path.Join(tempDir, "test.txt"))
	require.NoError(t, err)

	f.Close()

	require.Error(t, envx.IsDirEmpty(tempDir)(t.Context()))
}
