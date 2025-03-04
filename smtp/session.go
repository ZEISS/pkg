package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"sync"

	"github.com/katallaxie/pkg/group"
)

// Session ...
type Session struct {
	// Text is the textproto.Conn of the connection.
	Text *textproto.Conn
	// keep a reference to the connection so it can be used to create a TLS
	// connection later
	conn net.Conn
	// tls indicates whether the connection is already in TLS mode
	tls bool

	sync.Mutex
}

// NewSession...
func NewSession(conn net.Conn) *Session {
	s := new(Session)
	s.Text = textproto.NewConn(conn)
	s.conn = conn
	_, s.tls = conn.(*tls.Conn)

	return s
}

// Serve ...
func (s *Session) Serve() group.RunFunc {
	return func(ctx context.Context) {
		for {
			s.Text.PrintfLine("220 localhost ESMTP Service Ready")
			l, err := s.Text.ReadLine()
			if err != nil {
				return
			}

			fmt.Println(l)
		}
	}
}
