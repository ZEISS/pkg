package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMesage(t *testing.T) {
	t.Parallel()

	m, err := NewMessage()
	require.NoError(t, err)
	assert.NotNil(t, m)
	assert.NotEmpty(t, m.ID)
	assert.Equal(t, map[Header][]string{}, m.Headers)
}
