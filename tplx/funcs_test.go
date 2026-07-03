package tplx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTextFuncMap(t *testing.T) {
	t.Parallel()

	tfm := TxtFuncMap()
	require.NotNil(t, tfm)

	_, ok := tfm["hello"]
	require.True(t, ok)
}
