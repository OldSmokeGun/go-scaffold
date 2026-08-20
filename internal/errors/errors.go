package errors

import (
	stderrors "errors"
	"fmt"

	uerr "go-scaffold/pkg/errors"
)

// standard errors
var (
	ErrInternalError = New(50000, "internal error")

	ErrBadCall       = New(40000, "bad call")
	ErrValidateError = New(40001, "parameters validate error")

	ErrInvalidAuthorized = New(40100, "invalid authorized")

	ErrAccessDenied = New(40300, "access denied")

	ErrResourceNotFound = New(40400, "resource not found")

	ErrResourceConflict = New(40900, "resource conflict")

	ErrCallsTooFrequently = New(42900, "call too frequently")
)

// Error application internal error
type Error struct {
	code  int
	msg   string
	error error
}

// New returns an error that formats as the given text.
func New(code int, text string) *Error {
	return &Error{code, text, nil}
}

func (e *Error) Error() string {
	if e.error != nil {
		return e.error.Error()
	}
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.error
}

func (e *Error) Cause() error {
	return e.error
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) WithMsg(msg string) *Error {
	return &Error{e.code, msg, e.error}
}

func (e *Error) WithError(err error) *Error {
	return &Error{e.code, e.msg, err}
}

func (e *Error) Wrap(err error) error {
	if err == nil {
		return nil
	}
	return &Error{e.code, e.msg, uerr.WithStack(err, 1)}
}

func (e *Error) Errorf(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return &Error{e.code, msg, uerr.WithStack(stderrors.New(msg), 1)}
}
