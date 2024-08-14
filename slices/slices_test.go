package slices_test

import (
	"testing"

	"github.com/zeiss/pkg/slices"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		predicate func(v int) bool
		input     []int
		expected  bool
	}{
		{
			name:      "any element in slice",
			predicate: func(v int) bool { return v == 2 },
			input:     []int{1, 2, 3},
			expected:  true,
		},
		{
			name:      "any element not in slice",
			predicate: func(v int) bool { return v == 4 },
			input:     []int{1, 2, 3},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Any(tt.predicate, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		limit    int
		input    []int
		expected []int
	}{
		{
			name:     "limit slice",
			limit:    2,
			input:    []int{1, 2, 3},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Limit(tt.limit, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		from     int
		to       int
		expected []int
	}{
		{
			name:     "range slice",
			from:     1,
			to:       3,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Range(tt.from, tt.to)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fn       func(v int) int
		input    []int
		expected []int
	}{
		{
			name:     "map slice",
			fn:       func(v int) int { return v * 2 },
			input:    []int{1, 2, 3},
			expected: []int{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Map(tt.fn, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		el       int
		expected []int
	}{
		{
			name:     "pop from slice",
			input:    []int{1, 2, 3},
			el:       3,
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			el, actual := slices.Pop(tt.input...)
			assert.Equal(t, tt.el, el)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPush(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		el       int
		input    []int
		expected []int
	}{
		{
			name:     "push to slice",
			input:    []int{1, 2},
			el:       3,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Push(tt.el, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCut(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		start    int
		end      int
		input    []int
		expected []int
	}{
		{
			name:     "cut from slice",
			start:    1,
			end:      2,
			input:    []int{1, 2, 3},
			expected: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Cut(tt.start, tt.end, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		idx      int
		input    []int
		expected []int
	}{
		{
			name:     "delete from slice",
			idx:      1,
			input:    []int{1, 2, 3},
			expected: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Delete(tt.idx, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		idx      int
		el       int
		input    []int
		expected []int
	}{
		{
			name:     "insert into slice",
			idx:      1,
			el:       4,
			input:    []int{1, 2, 3},
			expected: []int{1, 4, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Insert(tt.el, tt.idx, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "filter slice",
			input:    []int{1, 2, 3},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Filter(func(el int) bool {
				return el < 3
			}, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestIn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		el       int
		input    []int
		expected bool
	}{
		{
			name:     "element in slice",
			el:       2,
			input:    []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "element not in slice",
			el:       4,
			input:    []int{1, 2, 3},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.In(tt.el, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		predicate func(v int) bool
		input     []int
		expected  int
	}{
		{
			name:      "element index in slice",
			predicate: func(v int) bool { return v == 2 },
			input:     []int{1, 2, 3},
			expected:  1,
		},
		{
			name:      "element index not in slice",
			predicate: func(v int) bool { return v == 4 },
			input:     []int{1, 2, 3},
			expected:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Index(tt.predicate, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "last element in slice",
			input:    []int{1, 2, 3},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Last(tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestUnique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     []int
		predicate func(v int) int
		expected  []int
	}{
		{
			name:      "unique elements in slice",
			predicate: func(v int) int { return v },
			input:     []int{1, 2, 2, 3},
			expected:  []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Unique(tt.predicate, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		size     int
		expected bool
	}{
		{
			name:     "empty slice",
			input:    []int{},
			size:     0,
			expected: true,
		},
		{
			name:     "slice with size",
			input:    []int{1, 2, 3},
			size:     3,
			expected: true,
		},
		{
			name:     "slice without size",
			input:    []int{1, 2, 3},
			size:     2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Size(tt.size, tt.input...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
