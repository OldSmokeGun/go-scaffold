package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"

	uerr "go-scaffold/pkg/errors"
)

// standard errors
var (
	ErrInternalError = New(50000, http.StatusInternalServerError, "internal error", "服务器出错")

	ErrBadCall       = New(40000, http.StatusBadRequest, "bad call", "客户端请求错误")
	ErrValidateError = New(40001, http.StatusBadRequest, "parameters validate error", "参数校验错误")

	ErrInvalidAuthorized = New(40100, http.StatusUnauthorized, "invalid authorized", "未经授权")

	ErrAccessDenied = New(40300, http.StatusForbidden, "access denied", "暂无权限")

	ErrResourceNotFound = New(40400, http.StatusNotFound, "resource not found", "资源不存在")

	ErrResourceConflict = New(40900, http.StatusConflict, "resource conflict", "资源冲突")

	ErrCallsTooFrequently = New(42900, http.StatusTooManyRequests, "call too frequently", "请求太频繁")
)

// Error application internal error
type Error struct {
	code       int
	httpStatus int
	hint       string
	msg        string
	error      error
}

func New(code, httpStatus int, msg, hint string) *Error {
	return &Error{code: code, httpStatus: httpStatus, msg: msg, hint: hint}
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

func (e *Error) HTTPStatus() int {
	return e.httpStatus
}

func (e *Error) HintMsg() string {
	return e.hint
}

func (e *Error) WithMsg(msg string) *Error {
	return e.spawn(msg, e.error)
}

func (e *Error) WithError(err error) *Error {
	return e.spawn(e.msg, err)
}

func (e *Error) Wrap(err error) error {
	if err == nil {
		return nil
	}
	return e.spawn(e.msg, uerr.WithStack(err, 4))
}

func (e *Error) Errorf(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return e.spawn(msg, uerr.WithStack(stderrors.New(msg), 4))
}

func (e *Error) spawn(msg string, err error) *Error {
	return &Error{code: e.code, httpStatus: e.httpStatus, hint: e.hint, msg: msg, error: err}
}
