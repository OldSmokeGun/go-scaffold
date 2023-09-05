package errors

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type HTTPError struct {
	*echo.HTTPError
}

func WrapHTTTPError(err *echo.HTTPError) *HTTPError {
	err.Internal = errors.WithStack(err.Internal)
	return &HTTPError{err}
}

func (e *HTTPError) SetMessage(message string) *HTTPError {
	e.Message = message
	return e
}

func (e *HTTPError) Unwrap() *echo.HTTPError {
	return e.HTTPError
}
