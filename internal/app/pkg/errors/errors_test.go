package errors

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"net/http"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "success", fields: fields{SuccessCode, SuccessCode.String()}, want: "code: 10000, message: OK"},
		{name: "server_error", fields: fields{ServerErrorCode, ServerErrorCode.String()}, want: "code: 10001, message: 服务器出错"},
		{name: "client_error", fields: fields{ClientErrorCode, ClientErrorCode.String()}, want: "code: 10002, message: 客户端请求错误"},
		{name: "validate_error", fields: fields{ValidateErrorCode, ValidateErrorCode.String()}, want: "code: 10003, message: 参数校验错误"},
		{name: "unauthorized", fields: fields{UnauthorizedCode, UnauthorizedCode.String()}, want: "code: 10004, message: 未经授权"},
		{name: "permission_denied", fields: fields{PermissionDeniedCode, PermissionDeniedCode.String()}, want: "code: 10005, message: 暂无权限"},
		{name: "resource_not_found", fields: fields{ResourceNotFoundCode, ResourceNotFoundCode.String()}, want: "code: 10006, message: 资源不存在"},
		{name: "too_many_request", fields: fields{TooManyRequestCode, TooManyRequestCode.String()}, want: "code: 10007, message: 请求过于频繁"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.fields.Code, tt.fields.Message)
			assert.Equal(t, tt.want, e.Error())
		})
	}

	t.Run("with_message", func(t *testing.T) {
		e := New(ServerErrorCode, ServerErrorCode.String(), WithMessage("test error"))
		want := "code: 10001, message: test error"
		assert.Equal(t, want, e.Error())
	})

	t.Run("with_metadata", func(t *testing.T) {
		e := New(SuccessCode, SuccessCode.String(), WithMetadata(map[string]string{"foo": "bar"}))
		want := "code: 10000, message: OK, metadata: {\"foo\": \"bar\"}"
		assert.Equal(t, want, e.Error())
	})
}

func TestStatusCode_HTTPStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       int
	}{
		{name: "http_status_ok", statusCode: SuccessCode.HTTPStatusCode(), want: http.StatusOK},
		{name: "http_status_internal_server_error", statusCode: ServerErrorCode.HTTPStatusCode(), want: http.StatusInternalServerError},
		{name: "http_status_bad_request", statusCode: ClientErrorCode.HTTPStatusCode(), want: http.StatusBadRequest},
		{name: "http_status_bad_request", statusCode: ValidateErrorCode.HTTPStatusCode(), want: http.StatusBadRequest},
		{name: "http_status_unauthorized", statusCode: UnauthorizedCode.HTTPStatusCode(), want: http.StatusUnauthorized},
		{name: "http_status_forbidden", statusCode: PermissionDeniedCode.HTTPStatusCode(), want: http.StatusForbidden},
		{name: "http_status_not_found", statusCode: ResourceNotFoundCode.HTTPStatusCode(), want: http.StatusNotFound},
		{name: "http_status_too_many_requests", statusCode: TooManyRequestCode.HTTPStatusCode(), want: http.StatusTooManyRequests},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.statusCode)
		})
	}
}

func TestError_GRPCStatus(t *testing.T) {
	t.Run("grpc_status", func(t *testing.T) {
		wantMetaData := map[string]string{
			"foo": "bar",
		}

		e := New(
			ServerErrorCode,
			ServerErrorCode.String(),
			WithMessage("test error"),
			WithMetadata(wantMetaData),
		)
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

		assert.Equal(t, wantMetaData, actualMetadata)
	})
}

func TestFromGRPCError(t *testing.T) {

}
