package goblerr

import "fmt"

// Error satisfies the built in error interface in addition to providing
// additional funcitonality for passing errors around the api
type Error struct {

	// Returns the error Code
	Code string `json:"message,omitempty"`

	// Returns the error message
	Message string `json:"code,omitempty"`

	// Returns the optional error details, nil if not set
	Detail interface{} `json:"detail,omitempty"`
}

// New returns a new error
func New(message, code string, detail interface{}) *Error {
	return &Error{Code: code, Message: message, Detail: detail}
}

func (b Error) Error() string {
	msg := b.Message

	if b.Detail != nil {
		msg = fmt.Sprintf("%s (%s)", msg, b.Detail)
	}

	return msg
}
