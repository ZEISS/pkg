package group

import (
	"context"
	"fmt"
	"sync"
)

// Group manages the lifetime of a set of goroutines from a common context.
// The first goroutine in the group to return will cause the context to be canceled,
// terminating the remaining goroutines.
type Group struct {
	// ctx is the context passed to all goroutines in the group.
	ctx    context.Context
	cancel context.CancelFunc
	done   sync.WaitGroup

	initOnce sync.Once

	errOnce sync.Once
	err     error
}

// Opt is an option for the routine group
type Opt func(*Group)

// WithContext uses the provided context for the group.
func WithContext(ctx context.Context) Opt {
	return func(g *Group) {
		g.ctx = ctx
	}
}

// New creates a new group.
func New(opts ...Opt) *Group {
	g := new(Group)

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// init initializes the group.
func (g *Group) init() {
	if g.ctx == nil {
		g.ctx = context.Background()
	}

	g.ctx, g.cancel = context.WithCancel(g.ctx)
}

// Add is adding a new goroutine to the group.
func (g *Group) Add(fn func(context.Context) error) {
	g.initOnce.Do(g.init)
	g.done.Add(1)

	go func() {
		defer g.done.Done()
		defer g.cancel()
		defer func() {
			if r := recover(); r != nil {
				g.errOnce.Do(func() {
					if err, ok := r.(error); ok {
						g.err = err
					} else {
						g.err = fmt.Errorf("panic: %v", r)
					}
				})
			}
		}()
		if err := fn(g.ctx); err != nil {
			g.errOnce.Do(func() { g.err = err })
		}
	}()
}

// Wait is a blocking call that waits for all goroutines to exit.
func (g *Group) Wait() error {
	g.done.Wait()
	g.errOnce.Do(func() {
		// noop, required to synchronise on the errOnce mutex.
	})

	return g.err
}
