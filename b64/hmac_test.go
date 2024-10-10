package b64_test

import (
	"testing"

	"github.com/zeiss/pkg/b64"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentHash(t *testing.T) {
	t.Parallel()

	content := []byte("test content")
	hash, err := b64.ContentHash(content)
	require.NoError(t, err)
	assert.Equal(t, "auinVVUgn9bEQVfArtgBbnY/9DWhnPGG92hjFAFD/3I=", hash)
}

func TestHmac256(t *testing.T) {
	t.Parallel()

	message := "test message"
	secret := "c2VjcmV0"

	hash, err := b64.Hmac256(message, secret)
	require.NoError(t, err)
	assert.Equal(t, "O86/Q8hdILum47a6J4rx0ro6sNV94nGwrTC4M+hRxaY=", hash)
}
