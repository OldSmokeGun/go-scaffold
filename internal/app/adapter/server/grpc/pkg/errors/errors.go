package errors

import (
	kerr "github.com/go-kratos/kratos/v2/errors"

	berr "go-scaffold/internal/errors"
)

var errMsg = map[error]string{
	berr.ErrInternalError:      "服务器出错",
	berr.ErrBadCall:            "客户端请求错误",
	berr.ErrValidateError:      "参数校验错误",
	berr.ErrInvalidAuthorized:  "未经授权",
	berr.ErrAccessDenied:       "暂无权限",
	berr.ErrResourceNotFound:   "资源不存在",
	berr.ErrCallsTooFrequently: "请求过于频繁",
}

func Message(err error) string {
	return errMsg[err]
}

// Wrap application internal error
func Wrap(err error) error {
	se, ok := err.(*berr.Error)
	if ok {
		return kerr.New(se.Code(), se.Label(), Message(se))
	}

	e := berr.ErrInternalError
	return kerr.New(e.Code(), e.Label(), err.Error())
}
