package duplex

import (
	"iter"
	"sync"

	"github.com/zeiss/pkg/slices"

	tea "charm.land/bubbletea/v2"
)

// Update is the type of the update of the store.
type Update any

// Action is the type of the action of the store.
type Action func() (Update, error)

// State is the type of the state of the store.
type State any

// Reduce is the type of the reducer of the store.
func Reduce[S State](curr S, reducers []Reducer[S], actions ...Action) iter.Seq[State] {
	return func(yield func(State) bool) {
		for _, action := range actions {
			for _, reducer := range reducers {
				curr = reducer(curr, action)
				if !yield(curr) {
					return
				}
			}
		}
	}
}

// Reducer is the type of the reducer of the store.
type Reducer[S State] func(curr S, update Update) S

// Duplexer is the interface for the duplex state machine.
type Duplexer[S State] interface {
	// Dispatch dispatches an event to the store.
	Dispatch(actions Action) tea.Cmd
}

// StateChangeMsg is a message that contains the.
type StateChangeMsg[S State] interface {
	// Prev gets the previous state of the store.
	Prev() S
	// Curr gets the current state of the store.
	Curr() S
}

// StateChangeError ...
type StateChangeError struct{}

type stateChange[S State] struct {
	prev S
	curr S
}

// NewStateChangeMsg creates a new state change.
func NewStateChangeMsg[S State](prev, curr S) StateChangeMsg[S] {
	return &stateChange[S]{
		prev: prev,
		curr: curr,
	}
}

// Curr gets the current state of the store.
func (s *stateChange[S]) Curr() S {
	return s.curr
}

// Prev gets the previous state of the store.
func (s *stateChange[S]) Prev() S {
	return s.prev
}

type duplexer[S State] struct {
	state    S
	reducers []Reducer[S]

	sync.RWMutex
}

// New creates a new store.
func New[S State](initialState S, reducers ...Reducer[S]) Duplexer[S] {
	s := new(duplexer[S])
	s.state = initialState
	s.reducers = slices.Append(reducers, s.reducers...)

	return s
}

// Dispatch dispatches an event to the store.
func (s *duplexer[S]) Dispatch(action Action) tea.Cmd {
	return func() tea.Msg {
		update, err := action()
		if err != nil {
			return &StateChangeError{}
		}

		prev := s.state

		for _, reducer := range s.reducers {
			s.state = reducer(s.state, update)
		}

		msg := NewStateChangeMsg(prev, s.state)

		return msg
	}
}
