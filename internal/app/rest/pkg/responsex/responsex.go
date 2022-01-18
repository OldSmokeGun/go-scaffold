package responsex

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
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

// Response 使用 application/json 类型返回 Body
func Response(ctx *gin.Context, c int, p *Body) {
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

var (
	// SuccessBody 成功响应 body
	SuccessBody = NewBody(SuccessCode, SuccessCode.String(), nil)
	// ServerErrorBody 服务器错误响应 body
	ServerErrorBody = NewBody(ServerErrorCode, ServerErrorCode.String(), nil)
	// ClientErrorBody 客户端错误响应 body
	ClientErrorBody = NewBody(ClientErrorCode, ClientErrorCode.String(), nil)
	// ValidateErrorBody 参数校验错误响应 body
	ValidateErrorBody = NewBody(ValidateErrorCode, ValidateErrorCode.String(), nil)
	// UnauthorizedBody 未经授权响应 body
	UnauthorizedBody = NewBody(UnauthorizedCode, UnauthorizedCode.String(), nil)
	// PermissionDeniedBody 暂无权限响应 body
	PermissionDeniedBody = NewBody(PermissionDeniedCode, PermissionDeniedCode.String(), nil)
	// ResourceNotFoundBody 资源不存在响应 body
	ResourceNotFoundBody = NewBody(ResourceNotFoundCode, ResourceNotFoundCode.String(), nil)
	// TooManyRequestBody 请求过于频繁响应 body
	TooManyRequestBody = NewBody(TooManyRequestCode, TooManyRequestCode.String(), nil)
)

// Success 成功响应
func Success(ctx *gin.Context, ops ...OptionFunc) {
	p := SuccessBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusOK, p)
}

// ServerError 服务器错误响应
func ServerError(ctx *gin.Context, ops ...OptionFunc) {
	p := ServerErrorBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusInternalServerError, p)
}

// ClientError 客户端错误响应
func ClientError(ctx *gin.Context, ops ...OptionFunc) {
	p := ClientErrorBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusBadRequest, p)
}

// ValidateError 参数校验错误响应
func ValidateError(ctx *gin.Context, ops ...OptionFunc) {
	p := ValidateErrorBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusBadRequest, p)
}

// Unauthorized 未经授权响应
func Unauthorized(ctx *gin.Context, ops ...OptionFunc) {
	p := UnauthorizedBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusUnauthorized, p)
}

// PermissionDenied 暂无权限响应
func PermissionDenied(ctx *gin.Context, ops ...OptionFunc) {
	p := PermissionDeniedBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusForbidden, p)
}

// ResourceNotFound 资源不存在响应
func ResourceNotFound(ctx *gin.Context, ops ...OptionFunc) {
	p := ResourceNotFoundBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusNotFound, p)
}

// TooManyRequest 请求过于频繁响应
func TooManyRequest(ctx *gin.Context, ops ...OptionFunc) {
	p := TooManyRequestBody
	for _, op := range ops {
		op(p)
	}
	Response(ctx, http.StatusTooManyRequests, p)
}
