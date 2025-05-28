package errors

import (
	"fmt"
)

type Error struct {
	Msg string
	Err error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
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

// IsNotFound checks if the error is a "not found" error
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "ent: user not found"
}
