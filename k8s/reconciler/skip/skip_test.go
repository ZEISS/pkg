package skip_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/pkg/k8s/reconciler/skip"
)

func TestSkipEnableSkip(t *testing.T) {
	ctx := context.Background()
	ctx = skip.EnableSkip(ctx)
	require.True(t, skip.Skip(ctx))
}
