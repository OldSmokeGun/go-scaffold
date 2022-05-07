package errors

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHTTPStatusCode(t *testing.T) {
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
