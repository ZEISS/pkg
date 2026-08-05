package sqlite_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/pkg/dbx/sqlite"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	path := filepath.Join(tempDir, ".builder", "builder.db")

	config := sqlite.NewConfig(path)
	require.NotNil(t, config)
}

func TestConfigDir(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	path := filepath.Join(tempDir, ".builder", "builder.db")
	config := sqlite.NewConfig(path)
	require.NotNil(t, config)

	dir, err := config.Dir()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(tempDir, ".builder"), dir)
}

func TestConfigMkdir(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	path := filepath.Join(tempDir, ".builder", "builder.db")
	config := sqlite.NewConfig(path)
	require.NotNil(t, config)

	err := config.Mkdir()
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(tempDir, ".builder"))
}
