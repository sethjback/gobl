package goblerr

import (
	"encoding/json"
	"fmt"
)

type baseError struct {
	Msg       string      `json:"message,omitempty"`
	ErrorCode string      `json:"code,omitempty"`
	Orig      error       `json:"origin,omitempty"`
	Det       interface{} `json:"detail,omitempty"`
}

func newBaseError(message, code string, origin error, detail interface{}) *baseError {
	return &baseError{message, code, origin, detail}
}

func (b baseError) Error() string {
	msg := b.Msg

	if b.Det != nil {
		msg = fmt.Sprintf("%s (%s)", msg, b.Detail())
	}

	if b.Orig != nil {
		msg = fmt.Sprintf("%s. caused by: %s", msg, b.Origin())
	}

	return msg
}

func (b baseError) Code() string {
	return b.ErrorCode
}

func (b baseError) Message() string {
	return b.Msg
}

func (b baseError) Origin() error {
	return b.Orig
}

func (b baseError) Detail() interface{} {
	return b.Det
}

func (b baseError) JSON() string {
	enc, e := json.Marshal(b)
	if e != nil {
		return e.Error()
	}

	return string(enc[:])
}
