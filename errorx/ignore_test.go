package errorx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/pkg/errorx"
)

func TestIgnore(t *testing.T) {
	val := errorx.Ignore(42, fmt.Errorf("an error occurred"))
	require.Equal(t, 42, val)
}

func TestNil(t *testing.T) {
	err := errorx.Nil(42)
	require.NoError(t, err)
}
