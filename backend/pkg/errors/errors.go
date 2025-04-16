package errors

import (
	"fmt"
)

type Error struct {
	Msg string
	Err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func WrapErr(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &Error{
		Msg: msg,
		Err: err,
	}
}

func New(msg string) error {
	return &Error{
		Msg: msg,
	}
}
