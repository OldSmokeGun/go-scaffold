package errors

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type HTTPError struct {
	*echo.HTTPError
}

func WrapHTTTPError[T error](err T) *HTTPError {
	target := &echo.HTTPError{}
	_ = errors.As(err, &target)
	return &HTTPError{target}
}

func (e *HTTPError) SetMessage(message string) *HTTPError {
	e.Message = message
	return e
}

func (e *HTTPError) Unwrap() error {
	return e.HTTPError
}
