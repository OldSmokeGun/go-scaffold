package errors

import (
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// Error 业务错误
type Error struct {
	// Code 状态码
	Code ErrorCode

	// Message 错误信息
	Message string

	// Metadata 元数据
	Metadata map[string]string

	cause error
}

func (e Error) Error() string {
	errStr := fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)

	if len(e.Metadata) > 0 {
		kvs := make([]string, 0, len(e.Message))

		for k, v := range e.Metadata {
			kvs = append(kvs, fmt.Sprintf("\"%s\": \"%s\"", k, v))
		}

		errStr = fmt.Sprintf("%s, metadata: {%s}", errStr, strings.Join(kvs, ", "))
	}

	return errStr
}

type Option func(e *Error)

func WithMessage(message string) Option {
	return func(e *Error) {
		e.Message = message
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(e *Error) {
		e.Metadata = metadata
	}
}

func New(code ErrorCode, message string, options ...Option) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}

	for _, option := range options {
		option(err)
	}

	return err
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
		ret := New(ErrorCode(gs.Code()), gs.Message())
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				ret.Metadata = d.Metadata
			}
		}
		return ret
	}
	return New(ServerErrorCode, err.Error())
}

// ServerError 服务器错误
func ServerError(options ...Option) *Error {
	return New(ServerErrorCode, ServerErrorCode.String(), options...)
}

// ClientError 客户端错误
func ClientError(options ...Option) *Error {
	return New(ClientErrorCode, ClientErrorCode.String(), options...)
}

// ValidateError 参数校验错误
func ValidateError(options ...Option) *Error {
	return New(ValidateErrorCode, ValidateErrorCode.String(), options...)
}

// Unauthorized 未认证
func Unauthorized(options ...Option) *Error {
	return New(UnauthorizedCode, UnauthorizedCode.String(), options...)
}

// PermissionDenied 权限拒绝错误
func PermissionDenied(options ...Option) *Error {
	return New(PermissionDeniedCode, PermissionDeniedCode.String(), options...)
}

// ResourceNotFound 资源不存在
func ResourceNotFound(options ...Option) *Error {
	return New(ResourceNotFoundCode, ResourceNotFoundCode.String(), options...)
}

// TooManyRequest 请求太过频繁
func TooManyRequest(options ...Option) *Error {
	return New(TooManyRequestCode, TooManyRequestCode.String(), options...)
}
