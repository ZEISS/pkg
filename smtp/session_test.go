package smtp_test

import (
	"context"
	"fmt"
	"net"
	s "net/smtp"
	"testing"
	"time"

	"github.com/katallaxie/pkg/group"
	"github.com/katallaxie/pkg/smtp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionServe(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	assert.NotNil(t, l)

	ctx, cancel := context.WithCancel(context.Background())
	g, _ := group.WithContext(ctx)

	ready := make(chan struct{})

	go func() {
		<-ready

		conn, err := net.DialTimeout("tcp", l.Addr().String(), 1*time.Second)
		require.NoError(t, err)

		c, err := s.NewClient(conn, "localhost")
		fmt.Println("hello")

		require.NoError(t, err)
		assert.NotNil(t, c)

		c.Close()
		cancel()
	}()

	close(ready)

	conn, err := l.Accept()
	require.NoError(t, err)
	assert.NotNil(t, conn)

	s := smtp.NewSession(conn)
	assert.NotNil(t, s)

	g.Run(s.Serve())
	g.Wait()
}
