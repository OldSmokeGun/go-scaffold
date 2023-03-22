package middleware

import (
	"net/http"

	berr "go-scaffold/internal/app/pkg/errors"
	uerr "go-scaffold/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

// ErrorHandler is HTTP error handler. It sends a JSON response
func ErrorHandler(debug bool, logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		logger.Error("handle request error", err)

		if ctx.Response().Committed {
			return
		}

		var (
			bc         int
			hintMsg    any
			statusCode int
		)

		ce := errors.Cause(err)

		switch er := ce.(type) {
		case *echo.HTTPError:
			statusCode = er.Code
			bc = berr.ErrServerError.Code()
			if c, ok := lo.Invert(errHttpStatusCode)[statusCode]; ok {
				bc = c
			}
			hintMsg = er.Message
		case *berr.Error:
			bc = er.Code()
			hintMsg = hintMessage(er.Label())
			statusCode = httpStatusCode(bc)
		default:
			de := berr.ErrServerError
			bc = de.Code()
			hintMsg = hintMessage(de.Label())
			statusCode = httpStatusCode(bc)
		}

		responseBody := NewDefaultBody().
			WithErrNo(bc).
			WithErrMsg(hintMsg)

		if debug {
			responseBody.WithErrMsg(err.Error())

			stack := uerr.ErrorStackTrace(err)
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
			logger.Error("send error response error", err)
		}
	}
}

var errHttpStatusCode = map[int]int{
	berr.ErrServerError.Code():      http.StatusInternalServerError,
	berr.ErrBadRequest.Code():       http.StatusBadRequest,
	berr.ErrValidateError.Code():    http.StatusBadRequest,
	berr.ErrUnauthorized.Code():     http.StatusUnauthorized,
	berr.ErrPermissionDenied.Code(): http.StatusForbidden,
	berr.ErrResourceNotFound.Code(): http.StatusNotFound,
	berr.ErrTooManyRequest.Code():   http.StatusTooManyRequests,
}

func httpStatusCode(c int) int {
	return errHttpStatusCode[c]
}

var errHintMsg = map[string]string{
	berr.ErrServerError.Label():      "服务器出错",
	berr.ErrBadRequest.Label():       "客户端请求错误",
	berr.ErrValidateError.Label():    "参数校验错误",
	berr.ErrUnauthorized.Label():     "未经授权",
	berr.ErrPermissionDenied.Label(): "暂无权限",
	berr.ErrResourceNotFound.Label(): "资源不存在",
	berr.ErrTooManyRequest.Label():   "请求过于频繁",
}

func hintMessage(l string) string {
	return errHintMsg[l]
}
