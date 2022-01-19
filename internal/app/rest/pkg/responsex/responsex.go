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
func NewSuccessBody() *Body { return NewBody(SuccessCode, SuccessCode.String(), nil) }

// NewServerErrorBody 服务器错误响应 body
func NewServerErrorBody() *Body { return NewBody(ServerErrorCode, ServerErrorCode.String(), nil) }

// NewClientErrorBody 客户端错误响应 body
func NewClientErrorBody() *Body { return NewBody(ClientErrorCode, ClientErrorCode.String(), nil) }

// NewValidateErrorBody 参数校验错误响应 body
func NewValidateErrorBody() *Body { return NewBody(ValidateErrorCode, ValidateErrorCode.String(), nil) }

// NewUnauthorizedBody 未经授权响应 body
func NewUnauthorizedBody() *Body { return NewBody(UnauthorizedCode, UnauthorizedCode.String(), nil) }

// NewPermissionDeniedBody 暂无权限响应 body
func NewPermissionDeniedBody() *Body {
	return NewBody(PermissionDeniedCode, PermissionDeniedCode.String(), nil)
}

// NewResourceNotFoundBody 资源不存在响应 body
func NewResourceNotFoundBody() *Body {
	return NewBody(ResourceNotFoundCode, ResourceNotFoundCode.String(), nil)
}

// NewTooManyRequestBody 请求过于频繁响应 body
func NewTooManyRequestBody() *Body {
	return NewBody(TooManyRequestCode, TooManyRequestCode.String(), nil)
}

// Success 成功响应
func Success(ctx *gin.Context, ops ...OptionFunc) {
	p := NewSuccessBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusOK, p)
}

// ServerError 服务器错误响应
func ServerError(ctx *gin.Context, ops ...OptionFunc) {
	p := NewServerErrorBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusInternalServerError, p)
}

// ClientError 客户端错误响应
func ClientError(ctx *gin.Context, ops ...OptionFunc) {
	p := NewClientErrorBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusBadRequest, p)
}

// ValidateError 参数校验错误响应
func ValidateError(ctx *gin.Context, ops ...OptionFunc) {
	p := NewValidateErrorBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusBadRequest, p)
}

// Unauthorized 未经授权响应
func Unauthorized(ctx *gin.Context, ops ...OptionFunc) {
	p := NewUnauthorizedBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusUnauthorized, p)
}

// PermissionDenied 暂无权限响应
func PermissionDenied(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPermissionDeniedBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusForbidden, p)
}

// ResourceNotFound 资源不存在响应
func ResourceNotFound(ctx *gin.Context, ops ...OptionFunc) {
	p := NewResourceNotFoundBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusNotFound, p)
}

// TooManyRequest 请求过于频繁响应
func TooManyRequest(ctx *gin.Context, ops ...OptionFunc) {
	p := NewTooManyRequestBody()
	for _, op := range ops {
		op(p)
	}
	ctx.JSON(http.StatusTooManyRequests, p)
}
