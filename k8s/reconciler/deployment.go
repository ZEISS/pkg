package reconciler

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// Event leverages go's 1.13 error wrapping.
type Event error

// EventIs reports whether any error in err's chain matches target.
var EventIs = errors.Is

// EventAs finds the first error in err's chain that matches target, and if so,
// sets target to that error value and returns true.
var EventAs = errors.As

// RecordFactory is a function that returns a new Event.
type RecordFactory func(eventtype, reason, messageFmt string, args ...interface{}) Event

// NewEvent returns a new Event.
func NewEvent(eventtype, reason, messageFmt string, args ...interface{}) Event {
	return &ReconcilerEvent{
		EventType: eventtype,
		Reason:    reason,
		Format:    messageFmt,
		Args:      args,
	}
}

// ReconcilerEvent ...
type ReconcilerEvent struct { //nolint:errname
	EventType string
	Reason    string
	Format    string
	Args      []interface{}
}

// make sure ReconcilerEvent implements error.
var _ error = (*ReconcilerEvent)(nil)

// Is ...
func (e *ReconcilerEvent) Is(target error) bool {
	var t *ReconcilerEvent
	if errors.As(target, &t) {
		if t != nil && t.EventType == e.EventType && t.Reason == e.Reason {
			return true
		}
		return false
	}
	// Allow for wrapped errors.
	err := fmt.Errorf(e.Format, e.Args...)
	return errors.Is(err, target)
}

// As ...
func (e *ReconcilerEvent) As(target interface{}) bool {
	err := fmt.Errorf(e.Format, e.Args...)
	return errors.As(err, target)
}

// Error ...
func (e *ReconcilerEvent) Error() string {
	return fmt.Errorf(e.Format, e.Args...).Error()
}

// Record ...
func (e *ReconcilerEvent) Record(obj runtime.Object, recorder record.EventRecorder) {
	recorder.Event(obj, e.EventType, e.Reason, e.Error())
}
