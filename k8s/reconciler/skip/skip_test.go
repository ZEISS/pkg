package skip_test

import (
	"testing"

	"github.com/zeiss/pkg/k8s/reconciler/skip"

	"github.com/stretchr/testify/require"
)

func TestSkipEnableSkip(t *testing.T) {
	ctx := t.Context()
	ctx = skip.EnableSkip(ctx)
	require.True(t, skip.Skip(ctx))
}
