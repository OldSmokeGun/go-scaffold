package errors

import (
	"errors"

	kerr "github.com/go-kratos/kratos/v2/errors"

	berr "go-scaffold/internal/errors"
)

var businessErrCodeMsgMap = map[int]string{
	berr.ErrInternalError.Code():      "服务器出错",
	berr.ErrBadCall.Code():            "客户端请求错误",
	berr.ErrValidateError.Code():      "参数校验错误",
	berr.ErrInvalidAuthorized.Code():  "未经授权",
	berr.ErrAccessDenied.Code():       "暂无权限",
	berr.ErrResourceNotFound.Code():   "资源不存在",
	berr.ErrResourceConflict.Code():   "资源冲突",
	berr.ErrCallsTooFrequently.Code(): "请求过于频繁",
}

var businessErrCodeReasonMap = map[int]string{
	berr.ErrInternalError.Code():      "INTERNAL_ERROR",
	berr.ErrBadCall.Code():            "BAD_CALL",
	berr.ErrValidateError.Code():      "VALIDATE_ERROR",
	berr.ErrInvalidAuthorized.Code():  "INVALID_AUTHORIZED",
	berr.ErrAccessDenied.Code():       "ACCESS_DENIED",
	berr.ErrResourceNotFound.Code():   "RESOURCE_NOT_FOUND",
	berr.ErrResourceConflict.Code():   "RESOURCE_CONFLICT",
	berr.ErrCallsTooFrequently.Code(): "CALLS_TOO_FREQUENTLY",
}

// Wrap application internal error
func Wrap(err error) error {
	var se *berr.Error
	if errors.As(err, &se) {
		return kerr.New(se.Code(), businessErrCodeReasonMap[se.Code()], businessErrCodeMsgMap[se.Code()])
	}

	e := berr.ErrInternalError
	return kerr.New(e.Code(), businessErrCodeReasonMap[e.Code()], err.Error())
}
