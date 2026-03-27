package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	berr "go-scaffold/internal/errors"
	perr "go-scaffold/pkg/errors"
)

// ErrorHandler is HTTP error handler. It sends a JSON response
func ErrorHandler(debug bool, logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		logger.Error("handle request error", slog.Any("error", err))

		var (
			httpErr *echo.HTTPError
			bErr    *berr.Error

			bc         int
			hintMsg    string
			statusCode int
		)

		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
			bc = berr.ErrInternalError.Code()
			if c, ok := httpStatusCodeBusinessErrCodeMap[statusCode]; ok {
				bc = c
			}
			hintMsg = fmt.Sprintf("%v", httpErr.Message)
			if une := httpErr.Unwrap(); une != nil {
				err = une
				var be *berr.Error
				if errors.As(une, &be) {
					if be.Unwrap() != nil {
						err = be.Unwrap()
					}
				}
			}
		} else if errors.As(err, &bErr) {
			bc = bErr.Code()
			hintMsg = bErr.Msg()
			statusCode = businessErrCodeHttpStatusCodeMap[bc]
			if bErr.Unwrap() != nil {
				err = bErr.Unwrap()
			}
		} else {
			de := berr.ErrInternalError
			bc = de.Code()
			hintMsg = businessErrCodeHintMsgMap[de.Code()]
			statusCode = businessErrCodeHttpStatusCodeMap[bc]
		}

		responseBody := NewDefaultBody().
			WithErrNo(bc).
			WithErrMsg(hintMsg)

		if debug {
			errMsg := err.Error()
			if hintMsg != "" {
				errMsg = fmt.Sprintf("%s: %s", hintMsg, err)
			}
			responseBody.WithErrMsg(errMsg)

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

var (
	businessErrCodeHttpStatusCodeMap = map[int]int{
		berr.ErrInternalError.Code():      http.StatusInternalServerError,
		berr.ErrBadCall.Code():            http.StatusBadRequest,
		berr.ErrValidateError.Code():      http.StatusBadRequest,
		berr.ErrInvalidAuthorized.Code():  http.StatusUnauthorized,
		berr.ErrAccessDenied.Code():       http.StatusForbidden,
		berr.ErrResourceNotFound.Code():   http.StatusNotFound,
		berr.ErrResourceConflict.Code():   http.StatusConflict,
		berr.ErrCallsTooFrequently.Code(): http.StatusTooManyRequests,
	}

	httpStatusCodeBusinessErrCodeMap = lo.Invert(businessErrCodeHttpStatusCodeMap)
)

var businessErrCodeHintMsgMap = map[int]string{
	berr.ErrInternalError.Code():      "服务器出错",
	berr.ErrBadCall.Code():            "客户端请求错误",
	berr.ErrValidateError.Code():      "参数校验错误",
	berr.ErrInvalidAuthorized.Code():  "未经授权",
	berr.ErrAccessDenied.Code():       "暂无权限",
	berr.ErrResourceNotFound.Code():   "资源不存在",
	berr.ErrResourceConflict.Code():   "资源冲突",
	berr.ErrCallsTooFrequently.Code(): "请求太频繁",
}
