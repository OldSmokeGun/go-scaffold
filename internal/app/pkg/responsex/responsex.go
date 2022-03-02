package responsex

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Body 响应格式
type Body struct {
	Code StatusCode  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NewBody 返回 *Body
func NewBody(code StatusCode, msg string, data interface{}) *Body {
	return &Body{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// OptionFunc Body 属性设置函数
type OptionFunc func(*Body)

// WithCode 设置 Body 的 Code
func WithCode(code StatusCode) OptionFunc {
	return func(p *Body) {
		p.Code = code
	}
}

// WithMsg 设置 Body 的 Msg
func WithMsg(msg string) OptionFunc {
	return func(p *Body) {
		p.Msg = msg
	}
}

// WithData 设置 Body 的 Data
func WithData(data interface{}) OptionFunc {
	return func(p *Body) {
		p.Data = data
	}
}

// NewSuccessBody 成功响应 body
func NewSuccessBody(ops ...OptionFunc) *Body {
	b := NewBody(SuccessCode, SuccessCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// NewServerErrorBody 服务器错误响应 body
func NewServerErrorBody(ops ...OptionFunc) *Body {
	b := NewBody(ServerErrorCode, ServerErrorCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// NewClientErrorBody 客户端错误响应 body
func NewClientErrorBody(ops ...OptionFunc) *Body {
	b := NewBody(ClientErrorCode, ClientErrorCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// NewValidateErrorBody 参数校验错误响应 body
func NewValidateErrorBody(ops ...OptionFunc) *Body {
	b := NewBody(ValidateErrorCode, ValidateErrorCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// NewUnauthorizedBody 未经授权响应 body
func NewUnauthorizedBody(ops ...OptionFunc) *Body {
	b := NewBody(UnauthorizedCode, UnauthorizedCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// NewPermissionDeniedBody 暂无权限响应 body
func NewPermissionDeniedBody(ops ...OptionFunc) *Body {
	b := NewBody(PermissionDeniedCode, PermissionDeniedCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// NewResourceNotFoundBody 资源不存在响应 body
func NewResourceNotFoundBody(ops ...OptionFunc) *Body {
	b := NewBody(ResourceNotFoundCode, ResourceNotFoundCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// NewTooManyRequestBody 请求过于频繁响应 body
func NewTooManyRequestBody(ops ...OptionFunc) *Body {
	b := NewBody(TooManyRequestCode, TooManyRequestCode.String(), nil)
	for _, op := range ops {
		op(b)
	}
	return b
}

// Success 成功响应
func Success(ctx *gin.Context, ops ...OptionFunc) {
	b := NewSuccessBody(ops...)
	ctx.JSON(http.StatusOK, b)
}

// ServerError 服务器错误响应
func ServerError(ctx *gin.Context, ops ...OptionFunc) {
	b := NewServerErrorBody(ops...)
	ctx.JSON(http.StatusInternalServerError, b)
}

// ClientError 客户端错误响应
func ClientError(ctx *gin.Context, ops ...OptionFunc) {
	b := NewClientErrorBody(ops...)
	ctx.JSON(http.StatusBadRequest, b)
}

// ValidateError 参数校验错误响应
func ValidateError(ctx *gin.Context, ops ...OptionFunc) {
	b := NewValidateErrorBody(ops...)
	ctx.JSON(http.StatusBadRequest, b)
}

// Unauthorized 未经授权响应
func Unauthorized(ctx *gin.Context, ops ...OptionFunc) {
	b := NewUnauthorizedBody(ops...)
	ctx.JSON(http.StatusUnauthorized, b)
}

// PermissionDenied 暂无权限响应
func PermissionDenied(ctx *gin.Context, ops ...OptionFunc) {
	b := NewPermissionDeniedBody(ops...)
	ctx.JSON(http.StatusForbidden, b)
}

// ResourceNotFound 资源不存在响应
func ResourceNotFound(ctx *gin.Context, ops ...OptionFunc) {
	b := NewResourceNotFoundBody(ops...)
	ctx.JSON(http.StatusNotFound, b)
}

// TooManyRequest 请求过于频繁响应
func TooManyRequest(ctx *gin.Context, ops ...OptionFunc) {
	b := NewTooManyRequestBody(ops...)
	ctx.JSON(http.StatusTooManyRequests, b)
}
