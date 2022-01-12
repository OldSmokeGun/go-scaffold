package responsex

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

// Payload 响应格式
type Payload struct {
	Code StatusCode  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NewPayload 返回 *Payload
func NewPayload(code StatusCode, msg string, data interface{}) *Payload {
	return &Payload{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// Response 使用 application/json 类型返回 Payload
func Response(ctx *gin.Context, c int, p *Payload) {
	if p.Data != nil {
		val := reflect.ValueOf(p.Data)
		switch val.Kind() {
		case reflect.Ptr:
			if val.IsNil() {
				p.Data = map[string]interface{}{}
			}
		case reflect.Slice:
			if val.IsNil() {
				p.Data = make([]interface{}, 0)
			}
		case reflect.Map:
			if val.IsNil() {
				p.Data = map[string]interface{}{}
			}
		case reflect.Chan:
			if val.IsNil() {
				p.Data = map[string]interface{}{}
			}
		case reflect.Func:
			if val.IsNil() {
				p.Data = map[string]interface{}{}
			}
		case reflect.Interface:
			if val.IsNil() {
				p.Data = map[string]interface{}{}
			}
		}
	} else {
		p.Data = map[string]interface{}{}
	}

	ctx.JSON(c, p)
}

// OptionFunc Payload 属性设置函数
type OptionFunc func(*Payload)

// WithCode 设置 Payload 的 Code
func WithCode(code StatusCode) OptionFunc {
	return func(p *Payload) {
		p.Code = code
	}
}

// WithMsg 设置 Payload 的 Msg
func WithMsg(msg string) OptionFunc {
	return func(p *Payload) {
		p.Msg = msg
	}
}

// WithData 设置 Payload 的 Data
func WithData(data interface{}) OptionFunc {
	return func(p *Payload) {
		p.Data = data
	}
}

// Success 成功响应
func Success(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(SuccessCode, SuccessCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusOK, p)
}

// ServerError 服务器错误响应
func ServerError(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(ServerErrorCode, ServerErrorCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusInternalServerError, p)
}

// ClientError 客户端错误响应
func ClientError(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(ClientErrorCode, ClientErrorCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusBadRequest, p)
}

// ValidateError 参数校验错误响应
func ValidateError(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(ValidateErrorCode, ValidateErrorCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusBadRequest, p)
}

// Unauthorized 登陆失效响应
func Unauthorized(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(UnauthorizedCode, UnauthorizedCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusUnauthorized, p)
}

// PermissionDenied 暂无权限响应
func PermissionDenied(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(PermissionDeniedCode, PermissionDeniedCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusForbidden, p)
}

// ResourceNotFound 资源不存在响应
func ResourceNotFound(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(ResourceNotFoundCode, ResourceNotFoundCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusNotFound, p)
}

// TooManyRequest 请求过于频繁响应
func TooManyRequest(ctx *gin.Context, ops ...OptionFunc) {
	p := NewPayload(TooManyRequestCode, TooManyRequestCode.String(), nil)
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusTooManyRequests, p)
}
