package smtp

import (
	"context"
	"net"
	"sync"

	"github.com/katallaxie/pkg/group"
	"github.com/katallaxie/pkg/server"
)

type mta struct {
	opts *Opts

	sync.Mutex
}

// MTA ...
type MTA interface {
	server.Listener
}

const (
	Network = "tcp"
)

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Addr string
}

// Configure ...
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// New ...
func New(opts ...Opt) *mta {
	options := new(Opts)
	options.Configure(opts...)

	return &mta{
		opts: options,
	}
}

// Start ...
func (m *mta) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		l, err := net.Listen(Network, m.opts.Addr)
		if err != nil {
			return err
		}

		g, _ := group.WithContext(ctx)

		for {
			c, err := l.Accept()
			if err != nil {
				return err
			}

			s := NewSession(c)
			g.Run(s.Serve())
		}
	}
}
