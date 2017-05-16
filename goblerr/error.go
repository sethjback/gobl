package goblerr

// Error satisfies the built in error interface in addition to providing
// additional funcitonality for passing errors around the api
type Error interface {
	error

	// Returns the error Code
	Code() string

	// Returns the error message
	Message() string

	// Returns the optional error details, nil if not set
	Detail() interface{}
}

// New returns a new error
func New(message, code string, detail interface{}) Error {
	return newBaseError(message, code, detail)
}
