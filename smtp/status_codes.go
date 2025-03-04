package smtp

// EnhancedMailSystemStatusCode ...
// https://datatracker.ietf.org/doc/html/rfc3463
type EnhancedMailSystemStatusCode [3]int

// StatusCode is a status code as defined in RFC 5321.
type StatusCode struct {
	replyCode    ReplyCode
	enhancedCode EnhancedMailSystemStatusCode
	message      string
}

// ReplyCode returns the reply code.
func (s *StatusCode) ReplyCode() int {
	return int(s.replyCode)
}

// EnhancedStatusCode returns the enhanced status code.
func (s *StatusCode) EnhancedStatusCode() EnhancedMailSystemStatusCode {
	return s.enhancedCode
}

// Message returns the message.
func (s *StatusCode) Message() string {
	return s.message
}

// Error returns error message.
func (s *StatusCode) Error() string {
	return s.message
}

// Error is an error as defined in RFC 5321.
type Error struct {
	statusCode *StatusCode
}

// Error returns error message.
func (e *Error) Error() string {
	return e.statusCode.message
}

// Temporary returns true if the error is temporary.
func (e *Error) Temporary() bool {
	return e.statusCode.replyCode/100 == 4
}

// NewStatusCode returns a new status code.
func NewStatusCode(replyCode ReplyCode, enhancedCode EnhancedMailSystemStatusCode, message string) *StatusCode {
	return &StatusCode{
		replyCode:    replyCode,
		enhancedCode: enhancedCode,
		message:      message,
	}
}

// Clone returns a clone of the status code.
func (s *StatusCode) Clone() *StatusCode {
	return &StatusCode{
		replyCode:    s.replyCode,
		enhancedCode: s.enhancedCode,
		message:      s.message,
	}
}

// Reset resets the status code.
func (s *StatusCode) Reset() {
	s.replyCode = 0
	s.enhancedCode = EnhancedMailSystemStatusCode{}
	s.message = ""
}

// ErrorFromStatus returns an error from a status code.
func ErrorFromStatus(s *StatusCode) error {
	return &Error{
		statusCode: s,
	}
}
