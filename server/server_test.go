package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	srv, ctx := WithContext(t.Context())
	assert.Implements(t, (*Server)(nil), srv)
	assert.NotNil(t, ctx)
	assert.NotNil(t, srv)
	assert.Nil(t, srv.sem)
}

func TestSetLimit(t *testing.T) {
	srv, ctx := WithContext(t.Context())
	assert.Implements(t, (*Server)(nil), srv)
	assert.NotNil(t, srv)
	assert.NotNil(t, ctx)

	srv.SetLimit(1)
	assert.NotNil(t, srv.sem)
}

func TestSetLimitZero(t *testing.T) {
	srv, ctx := WithContext(t.Context())
	assert.Implements(t, (*Server)(nil), srv)
	assert.NotNil(t, srv)
	assert.NotNil(t, ctx)

	srv.SetLimit(0)
	assert.NotNil(t, srv.sem)
}

func TestSetLimitNegative(t *testing.T) {
	srv, ctx := WithContext(t.Context())
	assert.Implements(t, (*Server)(nil), srv)
	assert.NotNil(t, srv)
	assert.NotNil(t, ctx)

	srv.SetLimit(-1)
	assert.Nil(t, srv.sem)
}

func TestUnimplemented(t *testing.T) {
	srv, ctx := WithContext(t.Context())
	assert.Implements(t, (*Server)(nil), srv)
	assert.NotNil(t, srv)
	assert.NotNil(t, ctx)

	l := &Unimplemented{}
	assert.Implements(t, (*Listener)(nil), l)

	srv.Listen(l, false)
	err := srv.Wait()
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnimplemented)
}

func TestNewError(t *testing.T) {
	err := NewServerError(ErrUnimplemented)
	assert.Implements(t, (*error)(nil), err)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnimplemented)
	require.Equal(t, "server: unimplemented", err.Error())
	require.NotNil(t, err.Unwrap())
	require.Equal(t, ErrUnimplemented, err.Unwrap())
}
