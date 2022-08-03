package casbin

import (
	"encoding/json"
	"errors"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/tests"
	"go-scaffold/internal/app/transport/http/pkg/response"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	jsonadapter "github.com/casbin/json-adapter/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	enforcer, err := setupEnforcer()
	if err != nil {
		t.Fatal(err)
	}

	cfg := New(enforcer, func(ctx *gin.Context) ([]any, error) {
		return nil, nil
	})

	assert.Equal(t, enforcer, cfg.Enforcer)
	assert.NotNil(t, cfg.RequestFunc)
	assert.Nil(t, cfg.ErrorResponseBody)
	assert.Nil(t, cfg.ValidateFailedResponseBody)
	assert.Nil(t, cfg.Logger)
}

func TestWithErrorResponseBody(t *testing.T) {
	var (
		cfg                   = &Casbin{}
		testErrorResponseBody = response.NewBody(int(errorsx.ServerErrorCode), errorsx.ServerErrorCode.String(), nil)
	)

	WithErrorResponseBody(testErrorResponseBody)(cfg)

	assert.Equal(t, testErrorResponseBody, cfg.ErrorResponseBody)
}

func TestWithLogger(t *testing.T) {
	ts, cleanup, err := tests.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	logger := log.NewHelper(ts.Logger)

	cfg := &Casbin{}

	WithLogger(logger)(cfg)

	assert.Equal(t, logger, cfg.Logger)
}

func TestWithValidateFailedResponseBody(t *testing.T) {
	var (
		cfg                            = &Casbin{}
		testValidateFailedResponseBody = response.NewBody(int(errorsx.UnauthorizedCode), errorsx.UnauthorizedCode.String(), nil)
	)

	WithValidateFailedResponseBody(testValidateFailedResponseBody)(cfg)

	assert.Equal(t, testValidateFailedResponseBody, cfg.ValidateFailedResponseBody)
}

func TestCasbin_Validate(t *testing.T) {
	ts, cleanup, err := tests.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	logger := log.NewHelper(ts.Logger)

	enforcer, err := setupEnforcer()
	if err != nil {
		t.Fatal(err)
	}

	options := []Option{
		WithLogger(logger),
		WithErrorResponseBody(response.NewBody(int(errorsx.ServerErrorCode), errorsx.ServerErrorCode.String(), nil)),
		WithValidateFailedResponseBody(response.NewBody(int(errorsx.PermissionDeniedCode), errorsx.PermissionDeniedCode.String(), nil)),
	}

	type args struct {
		enforcer    *casbin.Enforcer
		requestFunc requestFunc
	}

	type want struct {
		httpStatusCode int
		code           errorsx.ErrorCode
		msg            string
	}

	type uri struct {
		path   string
		method string
	}

	tts := []struct {
		name string
		uri  uri
		args args
		want want
	}{
		{
			name: "without_enforcer",
			uri:  uri{"/test", "GET"},
			args: args{nil, func(ctx *gin.Context) ([]any, error) {
				return nil, nil
			}},
			want: want{http.StatusInternalServerError, errorsx.ServerErrorCode, ErrNilEnforcer.Error()},
		},
		{
			name: "without_request_function",
			uri:  uri{"/test", "GET"},
			args: args{enforcer, nil},
			want: want{http.StatusInternalServerError, errorsx.ServerErrorCode, ErrNilRequestFunction.Error()},
		},
		{
			name: "get_casbin_request_parameters_error",
			uri:  uri{"/test", "GET"},
			args: args{enforcer, func(ctx *gin.Context) ([]any, error) {
				return nil, errors.New("test error")
			}},
			want: want{http.StatusInternalServerError, errorsx.ServerErrorCode, ErrGettingCasbinRequestParameters.Error()},
		},
		{
			name: "match_casbin_request_parameters_error",
			uri:  uri{"/test", "GET"},
			args: args{enforcer, func(ctx *gin.Context) ([]any, error) {
				return []any{"alice", ctx.Request.RequestURI, ctx.Request.Method, "match error"}, nil
			}},
			want: want{http.StatusInternalServerError, errorsx.ServerErrorCode, ErrMatchingCasbinRequestParameters.Error()},
		},
		{
			name: "validate_failed",
			uri:  uri{"/test", "GET"},
			args: args{enforcer, func(ctx *gin.Context) ([]any, error) {
				return []any{"bob", ctx.Request.RequestURI, ctx.Request.Method}, nil
			}},
			want: want{http.StatusForbidden, errorsx.PermissionDeniedCode, errorsx.PermissionDeniedCode.String()},
		},
		{
			name: "validate_success [/test GET]",
			uri:  uri{"/test", "GET"},
			args: args{enforcer, func(ctx *gin.Context) ([]any, error) {
				return []any{"alice", ctx.Request.RequestURI, ctx.Request.Method}, nil
			}},
			want: want{http.StatusOK, errorsx.SuccessCode, errorsx.SuccessCode.String()},
		},
		{
			name: "validate_success [/path/123 DELETE]",
			uri:  uri{"/path/123", "DELETE"},
			args: args{enforcer, func(ctx *gin.Context) ([]any, error) {
				return []any{"alice", ctx.Request.RequestURI, ctx.Request.Method}, nil
			}},
			want: want{http.StatusOK, errorsx.SuccessCode, errorsx.SuccessCode.String()},
		},
		{
			name: "validate_success [/parent/123/child/123 PUT]",
			uri:  uri{"/parent/123/child/123", "PUT"},
			args: args{enforcer, func(ctx *gin.Context) ([]any, error) {
				return []any{"alice", ctx.Request.RequestURI, ctx.Request.Method}, nil
			}},
			want: want{http.StatusOK, errorsx.SuccessCode, errorsx.SuccessCode.String()},
		},
		{
			name: "validate_success [/parent/123/child/456 PUT]",
			uri:  uri{"/parent/123/child/456", "PUT"},
			args: args{enforcer, func(ctx *gin.Context) ([]any, error) {
				return []any{"alice", ctx.Request.RequestURI, ctx.Request.Method}, nil
			}},
			want: want{http.StatusOK, errorsx.SuccessCode, errorsx.SuccessCode.String()},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.args.enforcer, tt.args.requestFunc, options...)

			router := setupRouter(c.Validate())

			r := httptest.NewRequest(tt.uri.method, tt.uri.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.want.httpStatusCode, w.Code)

			if w.Body.Bytes() != nil {
				body := new(response.Body)
				if err = json.Unmarshal(w.Body.Bytes(), body); err != nil {
					t.Fatal(err)
					return
				}

				assert.EqualValues(t, tt.want.code, body.Code)
				assert.Equal(t, tt.want.msg, body.Msg)
			}
		})
	}
}

func Test_handleResponse(t *testing.T) {
	var (
		testErrorResponseBody          = response.NewBody(int(errorsx.ServerErrorCode), errorsx.ServerErrorCode.String(), nil)
		testValidateFailedResponseBody = response.NewBody(int(errorsx.PermissionDeniedCode), errorsx.PermissionDeniedCode.String(), nil)
	)

	type want struct {
		httpStatusCode int
		code           errorsx.ErrorCode
		msg            string
	}

	tts := []struct {
		name        string
		handlerFunc gin.HandlerFunc
		want        want
	}{
		{
			"without_error_response_body",
			func(ctx *gin.Context) {
				handleResponse(ctx, http.StatusInternalServerError, nil, nil)
				return
			},
			want{http.StatusInternalServerError, 0, ""},
		},
		{
			"with_error_response_body_without_error",
			func(ctx *gin.Context) {
				handleResponse(ctx, http.StatusInternalServerError, testErrorResponseBody, nil)
				return
			},
			want{http.StatusInternalServerError, errorsx.ServerErrorCode, errorsx.ServerErrorCode.String()},
		},
		{
			"with_error_response_body_with_error",
			func(ctx *gin.Context) {
				handleResponse(ctx, http.StatusInternalServerError, testErrorResponseBody, errors.New("test error"))
				return
			},
			want{http.StatusInternalServerError, errorsx.ServerErrorCode, "test error"},
		},
		{
			"without_validate_failed_response_body",
			func(ctx *gin.Context) {
				handleResponse(ctx, http.StatusForbidden, nil, nil)
				return
			},
			want{http.StatusForbidden, 0, ""},
		},
		{
			"with_validate_failed_response_body",
			func(ctx *gin.Context) {
				handleResponse(ctx, http.StatusForbidden, testValidateFailedResponseBody, nil)
				return
			},
			want{http.StatusForbidden, errorsx.PermissionDeniedCode, errorsx.PermissionDeniedCode.String()},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.handlerFunc)

			r := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.want.httpStatusCode, w.Code)

			if w.Body.Bytes() != nil {
				body := new(response.Body)
				if err := json.Unmarshal(w.Body.Bytes(), body); err != nil {
					t.Fatal(err)
					return
				}

				assert.EqualValues(t, tt.want.code, body.Code)
				assert.Equal(t, tt.want.msg, body.Msg)
			}
		})
	}
}

func setupEnforcer() (*casbin.Enforcer, error) {
	modelConf := `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch3(r.obj, p.obj) && r.act == p.act`

	policy := []byte(`[
  {"PType":"p","V0":"admin","V1":"/test","V2":"GET"},
  {"PType":"p","V0":"admin","V1":"/path/{id}","V2":"DELETE"},
  {"PType":"p","V0":"admin","V1":"/parent/{id}/child/{id}","V2":"PUT"},
  {"PType":"g","V0":"alice","V1":"admin"},
  {"PType":"g","V0":"bob","V1":"guest"}
]`)

	mod, err := model.NewModelFromString(modelConf)
	if err != nil {
		return nil, err
	}

	adp := jsonadapter.NewAdapter(&policy)

	enforcer, err := casbin.NewEnforcer(mod, adp)
	if err != nil {
		return nil, err
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}

func setupRouter(mw gin.HandlerFunc) *gin.Engine {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/test", mw, func(c *gin.Context) { response.Success(c); return })
	r.DELETE("/path/:id", mw, func(c *gin.Context) { response.Success(c); return })
	r.PUT("/parent/:id/child/:id", mw, func(c *gin.Context) { response.Success(c); return })

	return r
}
