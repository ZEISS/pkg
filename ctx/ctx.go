package ctx

import (
	"context"
	"fmt"
	"reflect"

	"github.com/zeiss/pkg/cast"
)

// Key ...
type Key[Value any] struct {
	name *stringer[string]
	val  *Value
}

// New ...
func New[Value any](name string, v Value) Key[Value] {
	key := Key[Value]{}
	key.name = &stringer[string]{name}
	key.val = cast.Ptr(v)

	return key
}

// contextKey returns the context key to use.
func (k Key[Value]) contextKey() any {
	if k.name == nil {
		// Use the reflect.Type of the Value (implies key not created by New).
		return reflect.TypeFor[Value]()
	}

	return k.name
}

// WithValue  ...
func (k Key[Value]) WithValue(parent context.Context, val Value) context.Context {
	return context.WithValue(parent, k.contextKey(), stringer[Value]{val})
}

// ValueOk ...
func (k Key[Value]) ValueOk(ctx context.Context) (v Value, ok bool) {
	vv, ok := ctx.Value(k.contextKey()).(stringer[Value])
	if !ok && k.val != nil {
		vv.v = *k.val
	}

	return vv.v, ok
}

// Value ...
func (k Key[Value]) Value(ctx context.Context) (v Value) {
	v, _ = k.ValueOk(ctx)

	return v
}

// Has ...
func (k Key[Value]) Has(ctx context.Context) (ok bool) {
	_, ok = k.ValueOk(ctx)

	return ok
}

// String ...
func (k Key[Value]) String() string {
	if k.name == nil {
		return reflect.TypeFor[Value]().String()
	}

	return k.name.String()
}

// stringer implements [fmt.Stringer] on a generic T.
type stringer[T any] struct{ v T }

func (v stringer[T]) String() string { return fmt.Sprint(v.v) }
