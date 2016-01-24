package try

import "errors"

// Try wraps a funciton to handle mutiple retries
// https://github.com/matryer/try
type Try struct {
	MaxRetries int
}

var errMaxRetriesReached = errors.New("exceeded retry limit")

// New returns Try configured for max Retries
func New(maxRetries int) *Try {
	return &Try{maxRetries}
}

// Func represents functions that can be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func (t *Try) Do(fn Func) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > t.MaxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}

// IsMaxRetries checks whether the error is due to hitting the
// maximum number of retries or not.
func IsMaxRetries(err error) bool {
	return err == errMaxRetriesReached
}
