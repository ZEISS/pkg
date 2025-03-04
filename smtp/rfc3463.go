package smtp

import "fmt"

// EnhancedStatusCodeUnknown is the default value for the enhanced status code.
var EnhancedStatusCodeUnknown EnhancedMailSystemStatusCode = EnhancedMailSystemStatusCode{-1, -1, -1}

// EnhancedStatusCode is a data structure to contain enhanced
// mail system status codes from RFC 3463 (https://datatracker.ietf.org/doc/html/rfc3463).
type EnhancedStatusCode [3]int

// Class returns the class of the enhanced status code.
func (e EnhancedStatusCode) Class() int {
	return e[0]
}

// Subject returns the subject of the enhanced status code.
func (e EnhancedStatusCode) Subject() int {
	return e[1]
}

// Detail returns the detail of the enhanced status code.
func (e EnhancedStatusCode) Detail() int {
	return e[2]
}

// String returns the string representation of the enhanced status code.
func (e EnhancedMailSystemStatusCode) String() string {
	return fmt.Sprintf("%v.%v.%v", e[0], e[1], e[2])
}

const (
	// Signals a positive delivery action.
	EnhancedStatusCodeClassSuccess int = 2
	// Signals that there is a temporary failure in positively delivery action.
	EnhancedStatusCodeClassPersistentTransientFailure int = 4
	// Signals that there is a permanent failure in the delivery action.
	EnhancedStatusCodeClassPermanentFailure int = 5
)
