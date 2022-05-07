package response

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/pkg/errors"
	"net/http"
)

// Body 响应需实现的接口
type Body interface {
	WithCode(code int)
	WithMsg(msg string)
	WithData(data interface{})
}

// body 响应格式
type body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NewBody 返回 Body
func NewBody(code int, msg string, data interface{}) *body {
	return &body{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// WithCode 设置 body 的 Code
func (b *body) WithCode(code int) {
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

// Option body 属性设置函数
type Option func(*body)

// WithCode 设置 body 的 Code
func WithCode(code int) Option {
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

// Response HTTP 请求响应
func Response(ctx *gin.Context, httpStatusCode int, errorCode int, msg string, ops ...Option) {
	b := NewBody(errorCode, msg, nil)
	for _, op := range ops {
		op(b)
	}
	ctx.JSON(httpStatusCode, b)
}

// Success HTTP 成功响应
func Success(ctx *gin.Context, ops ...Option) {
	Response(ctx, http.StatusOK, int(errors.SuccessCode), errors.SuccessCode.String(), ops...)
}

// Error HTTP 错误响应
func Error(ctx *gin.Context, err error, ops ...Option) {
	ae, ok := err.(*errors.Error)
	if !ok {
		ae = errors.ServerError(errors.WithMessage(err.Error()))
	}
	Response(ctx, ae.Code.HTTPStatusCode(), int(ae.Code), ae.Message, ops...)
}
