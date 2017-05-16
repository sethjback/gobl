package goblerr

import (
	"fmt"
)

type baseError struct {
	Msg       string      `json:"message,omitempty"`
	ErrorCode string      `json:"code,omitempty"`
	Det       interface{} `json:"detail,omitempty"`
}

func newBaseError(message, code string, detail interface{}) *baseError {
	return &baseError{message, code, detail}
}

func (b baseError) Error() string {
	msg := b.Msg

	if b.Det != nil {
		msg = fmt.Sprintf("%s (%s)", msg, b.Detail())
	}

	return msg
}

func (b baseError) Code() string {
	return b.ErrorCode
}

func (b baseError) Message() string {
	return b.Msg
}

func (b baseError) Detail() interface{} {
	return b.Det
}
