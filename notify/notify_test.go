package notify

import (
	"context"
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

	ctx := context.Background()
	err := n.Notify(ctx, "subject", "message")
	require.NoError(t, err)
}
