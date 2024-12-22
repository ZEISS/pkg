package filex

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTempDir(t *testing.T) {
	f, fn, err := TempDir()
	require.NoError(t, err)
	defer fn()

	_, err = os.Stat(f.Name())
	os.IsNotExist(err)
	require.NoError(t, err)
}
