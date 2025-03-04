package smtp

import (
	"github.com/katallaxie/pkg/ulid"
)

// Header ...
type Header string

// Subject ...
const Subject Header = "Subject"

// From ...
const From Header = "From"

// To ...
const To Header = "To"

// Cc ...
const Cc Header = "Cc"

// Bcc ...
const Bcc Header = "Bcc"

// ReplyTo ...
const ReplyTo Header = "Reply-To"

// Date ...
const Date Header = "Date"

// Message ...
type Message struct {
	// ID ...
	ID string `json:"id" yaml:"id"`
	// Headers ...
	Headers map[Header][]string `json:"headers" yaml:"headers"`
}

// SetID ...
func (m *Message) SetID(id string) {
	m.ID = id
}

// NewMessage ...
func NewMessage() (*Message, error) {
	id, err := ulid.NewReverse()
	if err != nil {
		return nil, err
	}

	m := new(Message)
	m.SetID(id.String())
	m.Headers = map[Header][]string{}

	return m, nil
}
