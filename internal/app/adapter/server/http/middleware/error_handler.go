package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	berr "go-scaffold/internal/errors"
	perr "go-scaffold/pkg/errors"
)

// ErrorHandler is HTTP error handler. It sends a JSON response
func ErrorHandler(debug bool, logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		handleErr := err

		var (
			httpErr *echo.HTTPError
			bErr    *berr.Error

			bc         int
			hintMsg    string
			statusCode int
			cause      error
		)

		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
			bc = statusCode * 100
			hintMsg = fmt.Sprintf("%v", httpErr.Message)
			if une := httpErr.Unwrap(); une != nil {
				err = une
				cause = une
				var be *berr.Error
				if errors.As(une, &be) {
					if be.Unwrap() != nil {
						err = be.Unwrap()
						cause = be.Unwrap()
					}
				}
			}
		} else if errors.As(err, &bErr) {
			bc = bErr.Code()
			hintMsg = bErr.Msg()
			statusCode = bErr.HTTPStatus()
			cause = bErr.Cause()
			if bErr.Unwrap() != nil {
				err = bErr.Unwrap()
			}
		} else {
			de := berr.ErrInternalError
			bc = de.Code()
			hintMsg = de.Msg()
			statusCode = de.HTTPStatus()
		}

		if statusCode >= 500 {
			logger.Error("handle request error", slog.Any("error", handleErr))
		} else if statusCode >= 400 {
			logger.Warn("handle request error", slog.Any("error", handleErr))
		}

		responseBody := NewDefaultBody().
			WithErrNo(bc).
			WithErrMsg(hintMsg)

		if debug {
			if cause != nil {
				hintMsg = fmt.Sprintf("%s: %s", hintMsg, cause)
			}
			responseBody.WithErrMsg(hintMsg)

			stack := perr.ErrorStackTrace(err)
			if stack != nil {
				responseBody.WithStack(stack)
			}
		}

		if ctx.Request().Method == http.MethodHead { // Issue #608
			err = ctx.NoContent(statusCode)
		} else {
			err = ctx.JSON(statusCode, responseBody)
		}
		if err != nil {
			logger.Error("send error response error", slog.Any("error", err))
		}
	}
}
