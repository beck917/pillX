package utils

import (
	"errors"
)

type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewError(code int, msg string) *Error {
	err := &Error{}
	err.Code = code
	err.Err = errors.New(msg)

	//考虑自动记录日志

	return err
}
