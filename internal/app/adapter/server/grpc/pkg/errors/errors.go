package errors

import (
	berr "go-scaffold/internal/app/pkg/errors"

	kerr "github.com/go-kratos/kratos/v2/errors"
)

var errMsg = map[error]string{
	berr.ErrServerError:      "服务器出错",
	berr.ErrBadRequest:       "客户端请求错误",
	berr.ErrValidateError:    "参数校验错误",
	berr.ErrUnauthorized:     "未经授权",
	berr.ErrPermissionDenied: "暂无权限",
	berr.ErrResourceNotFound: "资源不存在",
	berr.ErrTooManyRequest:   "请求过于频繁",
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

	e := berr.ErrServerError
	return kerr.New(e.Code(), e.Label(), err.Error())
}
