package ctx

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"
)

// ErrCanceled is the returned when the CancelFunc returned by Merge is called.
var ErrCanceled = errors.New("context canceled")

// OneCtx is a context that can be used to wait for one of multiple contexts to be done.
type OneContext struct {
	ctx        context.Context
	ctxs       []context.Context
	done       chan struct{}
	err        error
	errMutex   sync.RWMutex
	cancelFunc context.CancelFunc
	cancelCtx  context.Context
}

var _ context.Context = (*OneContext)(nil)

// Merge merges the given contexts into a single context that will be done when any of the given contexts is done.
func Merge(ctx context.Context, ctxs ...context.Context) (context.Context, context.CancelFunc) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	o := &OneContext{
		done:       make(chan struct{}),
		ctx:        ctx,
		ctxs:       ctxs,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	go o.run()

	return o, cancelFunc
}

// Deadline returns the minimum deadline among all the contexts.
func (o *OneContext) Deadline() (time.Time, bool) {
	t := time.Time{}

	if deadline, ok := o.ctx.Deadline(); ok {
		t = deadline
	}

	for _, ctx := range o.ctxs {
		if deadline, ok := ctx.Deadline(); ok {
			if t.IsZero() || deadline.Before(t) {
				t = deadline
			}
		}
	}

	return t, !t.IsZero()
}

// Done returns a channel for cancellation.
func (o *OneContext) Done() <-chan struct{} {
	return o.done
}

// Err returns the first error raised by the contexts, otherwise a nil error.
func (o *OneContext) Err() error {
	o.errMutex.RLock()
	defer o.errMutex.RUnlock()

	return o.err
}

// Value returns the value associated with the key from one of the contexts.
func (o *OneContext) Value(key interface{}) interface{} {
	if value := o.ctx.Value(key); value != nil {
		return value
	}

	for _, ctx := range o.ctxs {
		if value := ctx.Value(key); value != nil {
			return value
		}
	}

	return nil
}

func (o *OneContext) run() {
	if len(o.ctxs) == 1 {
		o.runTwoContexts(o.ctx, o.ctxs[0])
		return
	}

	o.runMultipleContexts()
}

func (o *OneContext) runTwoContexts(ctx1, ctx2 context.Context) {
	select {
	case <-o.cancelCtx.Done():
		o.cancel(ErrCanceled)
	case <-ctx1.Done():
		o.cancel(o.ctx.Err())
	case <-ctx2.Done():
		o.cancel(o.ctxs[0].Err())
	}
}

func (o *OneContext) runMultipleContexts() {
	cases := make([]reflect.SelectCase, len(o.ctxs)+2)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(o.cancelCtx.Done())}
	cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(o.ctx.Done())}
	for i, ctx := range o.ctxs {
		cases[i+2] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())}
	}

	chosen, _, _ := reflect.Select(cases)
	switch chosen {
	case 0:
		o.cancel(ErrCanceled)
	case 1:
		o.cancel(o.ctx.Err())
	default:
		o.cancel(o.ctxs[chosen-2].Err())
	}
}

func (o *OneContext) cancel(err error) {
	o.errMutex.Lock()
	o.err = err
	o.errMutex.Unlock()

	close(o.done)
	o.cancelFunc()
}
