package notify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	n := New()
	require.NotNil(t, n)
}

func TestNotify(t *testing.T) {
	t.Parallel()

	n := New()
	require.NotNil(t, n)

	err := n.Notify(t.Context(), "subject", "message")
	require.NoError(t, err)
}
