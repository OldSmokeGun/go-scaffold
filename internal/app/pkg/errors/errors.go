package errors

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error 业务错误
type Error struct {
	// Code 错误状态码
	Code ErrorCode

	// Message 错误信息
	Message string

	// Metadata 元数据
	Metadata map[string]string

	cause error
}

// New 返回 *Error
func New(code ErrorCode, message string, metadata map[string]string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		Metadata: metadata,
	}
}

func (e *Error) Error() string {
	errStr := fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)

	if len(e.Metadata) > 0 {
		kvs := make([]string, 0, len(e.Message))

		for k, v := range e.Metadata {
			kvs = append(kvs, fmt.Sprintf("\"%s\": \"%s\"", k, v))
		}

		errStr = fmt.Sprintf("%s, metadata: {%s}", errStr, strings.Join(kvs, ", "))
	}

	if e.cause != nil {
		errStr = fmt.Sprintf("%s, cause: %v", errStr, e.cause)
	}

	return errStr
}

// Unwrap 兼容 Go 1.13 错误链
func (e *Error) Unwrap() error { return e.cause }

// WithCode 设置错误状态码
func (e *Error) WithCode(code ErrorCode) *Error {
	e.Code = code
	return e
}

// WithMessage 设置错误信息
func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

// WithMetadata 设置错误元数据
func (e *Error) WithMetadata(metadata map[string]string) *Error {
	e.Metadata = metadata
	return e
}

// WithCause 设置错误根本原因
func (e *Error) WithCause(cause error) *Error {
	e.cause = cause
	return e
}

// GRPCStatus 返回 GRPC 状态
func (e Error) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(e.Code), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Metadata: e.Metadata,
		})
	return s
}

// FromGRPCError 将 GRPC error 转换为 *Error
func FromGRPCError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := &Error{Code: ErrorCode(gs.Code()), Message: gs.Message()}
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				ret.Metadata = d.Metadata
			}
		}
		return ret
	}
	return &Error{Code: ServerErrorCode, Message: err.Error()}
}

// ServerError 服务器错误
func ServerError() *Error {
	return New(ServerErrorCode, ServerErrorCode.String(), nil)
}

// ClientError 客户端错误
func ClientError() *Error {
	return New(ClientErrorCode, ClientErrorCode.String(), nil)
}

// ValidateError 参数校验错误
func ValidateError() *Error {
	return New(ValidateErrorCode, ValidateErrorCode.String(), nil)
}

// Unauthorized 未认证
func Unauthorized() *Error {
	return New(UnauthorizedCode, UnauthorizedCode.String(), nil)
}

// PermissionDenied 权限拒绝错误
func PermissionDenied() *Error {
	return New(PermissionDeniedCode, PermissionDeniedCode.String(), nil)
}

// ResourceNotFound 资源不存在
func ResourceNotFound() *Error {
	return New(ResourceNotFoundCode, ResourceNotFoundCode.String(), nil)
}

// TooManyRequest 请求太过频繁
func TooManyRequest() *Error {
	return New(TooManyRequestCode, TooManyRequestCode.String(), nil)
}
