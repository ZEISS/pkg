package optional_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	optionals "github.com/zeiss/pkg/optional"
)

func TestOption_IsNone(t *testing.T) {
	a := optionals.Some(1)
	require.True(t, a.IsSome())

	b := optionals.None[any]()
	require.True(t, b.IsNone())
}

func TestOption_IsSome(t *testing.T) {
	a := optionals.Some(1)
	require.True(t, a.IsSome())

	b := optionals.None[any]()
	require.True(t, b.IsNone())
}

func TestOption_IsEmpty(t *testing.T) {
	a := optionals.Some[any](1)
	require.False(t, a.IsNone())
}

func TestOption_Unwrap(t *testing.T) {
	a := optionals.Some(1)
	require.Equal(t, 1, a.Unwrap())
}

func TestOption_UnwrapOrElse(t *testing.T) {
	a := optionals.Some(1)
	require.Equal(t, 1, a.UnwrapOrElse(func() int { return 2 }))

	c := optionals.Some[any](nil)
	require.Equal(t, 2, c.UnwrapOrElse(func() any { return 2 }))
}

func TestOption_And(t *testing.T) {
	a := optionals.Some(1)
	b := optionals.Some(2)
	require.Equal(t, 2, a.And(b).Unwrap())
}

func TestOption_Or(t *testing.T) {
	a := optionals.Some(1)
	b := optionals.Some(2)
	require.Equal(t, 1, a.Or(b).Unwrap())
}

func TestOption_Replace(t *testing.T) {
	a := optionals.Some(1)
	b := a.Replace(2)
	require.Equal(t, 1, b.Unwrap())
	require.Equal(t, 2, a.Unwrap())
}

func TestOption_Map(t *testing.T) {
	a := optionals.Some(1)
	b := a.Map(func(i int) int { return i + 1 })
	require.Equal(t, 2, b.Unwrap())
	require.Equal(t, 1, a.Unwrap())
}

func TestOption_MapOr(t *testing.T) {
	a := optionals.Some(1)
	require.Equal(t, 2, a.MapOr(3, func(i int) int { return i + 1 }))

	b := optionals.Some[interface{}](nil)
	require.Equal(t, 3, b.MapOr(3, func(i any) any { return i }))
}

func TestOption_MapOrElse(t *testing.T) {
	a := optionals.Some(1)
	require.Equal(t, 2, a.MapOrElse(func() int { return 3 }, func(i int) int { return i + 1 }))

	b := optionals.Some[any](nil)
	require.Equal(t, 3, b.MapOrElse(func() any { return 3 }, func(i any) any { return i }))
}
