package responsex

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Body 响应需实现的接口
type Body interface {
	WithCode(code StatusCode)
	WithMsg(msg string)
	WithData(data interface{})
}

// body 响应格式
type body struct {
	Code StatusCode `json:"code"`
	Msg  string     `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NewBody 返回 Body
func NewBody(code StatusCode, msg string, data interface{}) *body {
	return &body{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// WithCode 设置 body 的 Code
func (b *body) WithCode(code StatusCode) {
	b.Code = code
}

// WithMsg 设置 body 的 Msg
func (b *body) WithMsg(msg string) {
	b.Msg = msg
}

// WithData 设置 body 的 Data
func (b *body) WithData(data interface{}) {
	b.Data = data
}

// NewSuccessBody 成功响应 body
func NewSuccessBody() *body {
	b := NewBody(SuccessCode, SuccessCode.String(), nil)
	return b
}

// NewServerErrorBody 服务器错误响应 body
func NewServerErrorBody() *body {
	b := NewBody(ServerErrorCode, ServerErrorCode.String(), nil)
	return b
}

// NewClientErrorBody 客户端错误响应 body
func NewClientErrorBody() *body {
	b := NewBody(ClientErrorCode, ClientErrorCode.String(), nil)
	return b
}

// NewValidateErrorBody 参数校验错误响应 body
func NewValidateErrorBody() *body {
	b := NewBody(ValidateErrorCode, ValidateErrorCode.String(), nil)
	return b
}

// NewUnauthorizedBody 未经授权响应 body
func NewUnauthorizedBody() *body {
	b := NewBody(UnauthorizedCode, UnauthorizedCode.String(), nil)
	return b
}

// NewPermissionDeniedBody 暂无权限响应 body
func NewPermissionDeniedBody() *body {
	b := NewBody(PermissionDeniedCode, PermissionDeniedCode.String(), nil)
	return b
}

// NewResourceNotFoundBody 资源不存在响应 body
func NewResourceNotFoundBody() *body {
	b := NewBody(ResourceNotFoundCode, ResourceNotFoundCode.String(), nil)
	return b
}

// NewTooManyRequestBody 请求过于频繁响应 body
func NewTooManyRequestBody() *body {
	b := NewBody(TooManyRequestCode, TooManyRequestCode.String(), nil)
	return b
}

// Option body 属性设置函数
type Option func(*body)

// WithCode 设置 body 的 Code
func WithCode(code StatusCode) Option {
	return func(p *body) {
		p.Code = code
	}
}

// WithMsg 设置 body 的 Msg
func WithMsg(msg string) Option {
	return func(p *body) {
		p.Msg = msg
	}
}

// WithData 设置 body 的 Data
func WithData(data interface{}) Option {
	return func(p *body) {
		p.Data = data
	}
}

// Success 成功响应
func Success(ctx *gin.Context, ops ...Option) {
	b := NewSuccessBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusOK, b)
}

// ServerError 服务器错误响应
func ServerError(ctx *gin.Context, ops ...Option) {
	b := NewServerErrorBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusInternalServerError, b)
}

// ClientError 客户端错误响应
func ClientError(ctx *gin.Context, ops ...Option) {
	b := NewClientErrorBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusBadRequest, b)
}

// ValidateError 参数校验错误响应
func ValidateError(ctx *gin.Context, ops ...Option) {
	b := NewValidateErrorBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusBadRequest, b)
}

// Unauthorized 未经授权响应
func Unauthorized(ctx *gin.Context, ops ...Option) {
	b := NewUnauthorizedBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusUnauthorized, b)
}

// PermissionDenied 暂无权限响应
func PermissionDenied(ctx *gin.Context, ops ...Option) {
	b := NewPermissionDeniedBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusForbidden, b)
}

// ResourceNotFound 资源不存在响应
func ResourceNotFound(ctx *gin.Context, ops ...Option) {
	b := NewResourceNotFoundBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusNotFound, b)
}

// TooManyRequest 请求过于频繁响应
func TooManyRequest(ctx *gin.Context, ops ...Option) {
	b := NewTooManyRequestBody()
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(http.StatusTooManyRequests, b)
}
