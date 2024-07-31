package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ErrUnimplemented is returned when a listener is not implemented.
var ErrUnimplemented = errors.New("server: unimplemented")

type token struct{}

// ReadyFunc is the function that is called by Listener to signal
// that it is ready and the next Listener can be called.
type ReadyFunc func()

// RunFunc is a function that is called to attach more routines
// to the server.
type RunFunc func(func() error)

// ServerError ...
type ServerError struct {
	Err error
}

// Error implements the error interface.
func (s *ServerError) Error() string { return fmt.Sprintf("server: %s", s.Err) }

// Unwrap ...
func (s *ServerError) Unwrap() error { return s.Err }

// NewServerError returns a new error.
func NewServerError(err error) *ServerError {
	return &ServerError{Err: err}
}

// Server is the interface to be implemented
// to run the server.
//
//	s, ctx := WithContext(context.Background())
//	s.Listen(listener, false)
//
//	if err := s.Wait(); err != nil {
//		panic(err)
//	}
type Server interface {
	// Run is running a new go routine
	Listen(listener Listener, ready bool)

	// Waits for the server to fail,
	// or gracefully shutdown if context is canceled
	Wait() error

	// SetLimit ...
	SetLimit(n int)
}

// Unimplemented is the default implementation.
type Unimplemented struct{}

// Start ...
func (s *Unimplemented) Start(context.Context, ReadyFunc, RunFunc) func() error {
	return func() error {
		return ErrUnimplemented
	}
}

// Listener is the interface to a listener,
// so starting and shutdown of a listener,
// or any routine.
type Listener interface {
	// Start is being called on the listener
	Start(context.Context, ReadyFunc, RunFunc) func() error
}

type listeners map[Listener]bool

// server holds the instance info of the server
type server struct {
	ctx    context.Context
	cancel context.CancelFunc

	wg      sync.WaitGroup
	errOnce sync.Once
	err     error

	sem chan token

	listeners map[Listener]bool

	ready chan bool
	sys   chan os.Signal
}

// WithContext is creating a new server with a context.
func WithContext(ctx context.Context) (*server, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	// new server
	s := newServer(ctx)
	s.cancel = cancel
	s.ctx = ctx

	return s, ctx
}

func newServer(ctx context.Context) *server {
	s := new(server)

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.ctx = ctx

	s.listeners = make(listeners)
	s.ready = make(chan bool, 1)
	s.sys = make(chan os.Signal, 1)

	return s
}

// Listen is adding a listener to the server.
func (s *server) Listen(listener Listener, ready bool) {
	s.listeners[listener] = ready
}

// Wait is waiting for the server to shutdown or fail.
// The returned error is the first error that occurred from the listeners.
//
//nolint:gocyclo
func (s *server) Wait() error {
	// create ticker for interrupt signals
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	signal.Notify(s.sys, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Reset(syscall.SIGINT, syscall.SIGTERM)

OUTTER:
	//
	for l, ready := range s.listeners {
		readyFunc := func() {
			r := ready

			var readyOnce sync.Once
			readyOnce.Do(func() {
				if r {
					s.ready <- true
				}
			})
		}

		goFn := func(f func() error) { s.run(f) }

		// schedule to routines
		s.run(l.Start(s.ctx, readyFunc, goFn))

		// this blocks until ready is called
		if ready {
			select {
			case <-s.ready:
				continue OUTTER
			case <-s.sys:
				s.cancel()
				break OUTTER
			case <-s.ctx.Done():
				break OUTTER
			}
		}
	}

	for {
		select {
		case <-ticker.C:
		case <-s.sys:
			// if there is sys interrupt
			// cancel the context of the routines
			s.cancel()
		case <-s.ctx.Done():
			if s.err != nil {
				return s.err
			}

			return nil
		}
	}
}

// SetLimit limits the number of active listeners in this server
func (s *server) SetLimit(n int) {
	if n < 0 {
		s.sem = nil
		return
	}

	if len(s.sem) != 0 {
		panic(fmt.Errorf("server: modify limit while %v listeners run", len(s.sem)))
	}

	s.sem = make(chan token, n)
}

func (s *server) run(f func() error) {
	if s.sem != nil {
		s.sem <- token{}
	}

	s.wg.Add(1)

	fn := func() {
		defer s.done()

		if err := f(); err != nil {
			s.errOnce.Do(func() {
				s.err = err
				if s.cancel != nil {
					s.cancel()
				}
			})
		}
	}

	go fn()
}

func (s *server) done() {
	if s.sem != nil {
		<-s.sem
	}

	s.wg.Done()
}
