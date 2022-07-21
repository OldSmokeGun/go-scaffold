package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name  string
		error *Error
		want  string
	}{
		{name: "success", error: New(SuccessCode, SuccessCode.String(), nil), want: "code: 10000, message: OK"},
		{name: "server_error", error: ServerError(), want: "code: 10001, message: 服务器出错"},
		{name: "client_error", error: ClientError(), want: "code: 10002, message: 客户端请求错误"},
		{name: "validate_error", error: ValidateError(), want: "code: 10003, message: 参数校验错误"},
		{name: "unauthorized", error: Unauthorized(), want: "code: 10004, message: 未经授权"},
		{name: "permission_denied", error: PermissionDenied(), want: "code: 10005, message: 暂无权限"},
		{name: "resource_not_found", error: ResourceNotFound(), want: "code: 10006, message: 资源不存在"},
		{name: "too_many_request", error: TooManyRequest(), want: "code: 10007, message: 请求过于频繁"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.error.Code, tt.error.Message, nil)
			assert.Equal(t, tt.want, e.Error())
		})
	}

	t.Run("with_code", func(t *testing.T) {
		e := New(0, "test error", nil).WithCode(ServerErrorCode)
		want := "code: 10001, message: test error"
		assert.Equal(t, want, e.Error())
	})

	t.Run("with_message", func(t *testing.T) {
		e := New(ServerErrorCode, ServerErrorCode.String(), nil).WithMessage("test error")
		want := "code: 10001, message: test error"
		assert.Equal(t, want, e.Error())
	})

	t.Run("with_metadata", func(t *testing.T) {
		e := New(SuccessCode, SuccessCode.String(), nil).WithMetadata(map[string]string{"foo": "bar"})
		want := "code: 10000, message: OK, metadata: {\"foo\": \"bar\"}"
		assert.Equal(t, want, e.Error())
	})

	t.Run("with_cause", func(t *testing.T) {
		e := New(ServerErrorCode, ServerErrorCode.String(), nil).WithMessage("test error").WithCause(errors.New("test cause error"))
		want := "code: 10001, message: test error, cause: test cause error"
		assert.Equal(t, want, e.Error())
	})

	t.Run("unwrap_error", func(t *testing.T) {
		testCause := errors.New("test cause error")
		e := New(ServerErrorCode, ServerErrorCode.String(), nil).WithCause(testCause)
		assert.Equal(t, testCause, errors.Unwrap(e))
	})
}

func TestGRPCStatus(t *testing.T) {
	t.Run("grpc_status", func(t *testing.T) {
		metaData := map[string]string{
			"foo": "bar",
		}

		e := New(ServerErrorCode, ServerErrorCode.String(), nil).
			WithMessage("test error").
			WithMetadata(metaData)
		gs := e.GRPCStatus()

		assert.Equal(t, 10001, int(gs.Code()))
		assert.Equal(t, "test error", gs.Message())

		var actualMetadata map[string]string
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				actualMetadata = d.Metadata
			}
		}

		assert.Equal(t, metaData, actualMetadata)
	})
}

func TestFromGRPCError(t *testing.T) {
	t.Run("error_is_nil", func(t *testing.T) {
		var err error
		e := FromGRPCError(err)
		assert.Nil(t, e)
	})

	t.Run("error", func(t *testing.T) {
		metadata := map[string]string{"foo": "bar"}
		err := New(ServerErrorCode, ServerErrorCode.String(), nil).WithMetadata(metadata)
		e := FromGRPCError(err)
		assert.Equal(t, ServerErrorCode, e.Code)
		assert.Equal(t, ServerErrorCode.String(), e.Message)
		assert.Equal(t, metadata, e.Metadata)
	})

	t.Run("grpc_error", func(t *testing.T) {
		metadata := map[string]string{"foo": "bar"}
		err := New(ServerErrorCode, ServerErrorCode.String(), nil).WithMetadata(metadata)
		e := FromGRPCError(err.GRPCStatus().Err())
		assert.Equal(t, ServerErrorCode, e.Code)
		assert.Equal(t, ServerErrorCode.String(), e.Message)
		assert.Equal(t, metadata, e.Metadata)
	})

	t.Run("other_error", func(t *testing.T) {
		var err = errors.New("test error")
		e := FromGRPCError(err)
		assert.Equal(t, ServerErrorCode, e.Code)
		assert.Equal(t, "test error", e.Message)
	})
}
