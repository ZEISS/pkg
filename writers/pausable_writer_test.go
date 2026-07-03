package writers_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/pkg/writers"
)

func TestPausableWriter(t *testing.T) {
	var buf bytes.Buffer

	pw := writers.NewPausableWriter(&buf)
	require.NotNil(t, pw)

	n, err := pw.Write([]byte("Hello, world!"))
	require.NoError(t, err)
	require.Equal(t, 13, n)

	pw.Pause()

	n, err = pw.Write([]byte("This should not be written"))
	require.NoError(t, err)
	require.Equal(t, 26, n)

	pw.Resume()

	n, err = pw.Write([]byte("This should be written"))
	require.NoError(t, err)
	require.Equal(t, 22, n)

	expected := "Hello, world!This should be written"
	require.Equal(t, expected, buf.String())
}
