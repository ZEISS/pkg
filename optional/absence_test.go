package optional_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeiss/pkg/optional"
)

func TestAndPresenceCombinations(t *testing.T) {
	tests := []struct {
		name    string
		left    optional.Option[int]
		right   optional.Option[int]
		present bool
		value   int
	}{
		{"some and some", optional.Some(1), optional.Some(2), true, 2},
		{"some and none", optional.Some(1), optional.None[int](), false, 0},
		{"none and some", optional.None[int](), optional.Some(2), false, 0},
		{"none and none", optional.None[int](), optional.None[int](), false, 0},
		{"some and zero", optional.Some(1), optional.Some(0), true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.left.And(tt.right)
			assert.Equal(t, tt.present, result.IsSome())
			if tt.present {
				assert.Equal(t, tt.value, result.Unwrap())
			}
		})
	}
}

func TestFilterAbsentDoesNotCallPredicate(t *testing.T) {
	for _, option := range []optional.Option[*int]{optional.None[*int](), optional.Some((*int)(nil))} {
		calls := 0
		result := option.Filter(func(value *int) bool {
			calls++
			return value != nil
		})
		assert.True(t, result.IsNone())
		assert.Zero(t, calls)
	}
}

func TestFilterPresentCallsPredicateOnce(t *testing.T) {
	for _, keep := range []bool{false, true} {
		calls := 0
		result := optional.Some(0).Filter(func(value int) bool {
			calls++
			assert.Zero(t, value)
			return keep
		})
		assert.Equal(t, keep, result.IsSome())
		assert.Equal(t, 1, calls)
		if keep {
			assert.Zero(t, result.Unwrap())
		}
	}
}

func TestReplaceAbsentAndPresent(t *testing.T) {
	tests := []struct {
		name            string
		option          optional.Option[int]
		replacement     int
		previousPresent bool
		previousValue   int
	}{
		{"none to value", optional.None[int](), 2, false, 0},
		{"none to zero", optional.None[int](), 0, false, 0},
		{"some to value", optional.Some(1), 2, true, 1},
		{"some zero to value", optional.Some(0), 2, true, 0},
		{"wrapped none to value", optional.Wrap(7, true), 2, false, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			previous := tt.option.Replace(tt.replacement)
			assert.Equal(t, tt.previousPresent, previous.IsSome())
			assert.Equal(t, tt.previousValue, previous.Value())
			assert.True(t, tt.option.IsSome())
			assert.Equal(t, tt.replacement, tt.option.UnwrapOr(-1))
			tt.option.Replace(3)
			assert.Equal(t, tt.previousPresent, previous.IsSome())
			assert.Equal(t, tt.previousValue, previous.Value())
		})
	}
}
