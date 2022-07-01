package jwt

import (
	"bou.ke/monkey"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/tests"
	"go-scaffold/internal/app/transport/http/pkg/response"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := New("test key")

	assert.Equal(t, "test key", cfg.Key)
	assert.Equal(t, defaultHeaderName, cfg.HeaderName)
	assert.Equal(t, defaultHeaderPrefix, cfg.HeaderPrefix)
	assert.Nil(t, cfg.ErrorResponseBody)
	assert.Nil(t, cfg.ValidateFailedResponseBody)
	assert.Nil(t, cfg.Logger)
	assert.NotNil(t, cfg.PostFunc)
	assert.Equal(t, "", cfg.Raw)
}

func TestWithHeaderName(t *testing.T) {
	var (
		cfg            = &JWT{}
		testHeaderName = "TestHeaderName"
	)

	WithHeaderName(testHeaderName)(cfg)

	assert.Equal(t, testHeaderName, cfg.HeaderName)
}

func TestWithHeaderPrefix(t *testing.T) {
	var (
		cfg              = &JWT{}
		testHeaderPrefix = "Test"
	)

	WithHeaderPrefix(testHeaderPrefix)(cfg)

	assert.Equal(t, testHeaderPrefix, cfg.HeaderPrefix)
}

func TestWithErrorResponseBody(t *testing.T) {
	var (
		cfg                   = &JWT{}
		testErrorResponseBody = response.NewBody(int(errorsx.ServerErrorCode), errorsx.ServerErrorCode.String(), nil)
	)

	WithErrorResponseBody(testErrorResponseBody)(cfg)

	assert.Equal(t, testErrorResponseBody, cfg.ErrorResponseBody)
}

func TestWithValidateFailedResponseBody(t *testing.T) {
	var (
		cfg                            = &JWT{}
		testValidateFailedResponseBody = response.NewBody(int(errorsx.UnauthorizedCode), errorsx.UnauthorizedCode.String(), nil)
	)

	WithValidateFailedResponseBody(testValidateFailedResponseBody)(cfg)

	assert.Equal(t, testValidateFailedResponseBody, cfg.ValidateFailedResponseBody)
}

func TestWithLogger(t *testing.T) {
	ts, cleanup, err := tests.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	logger := log.NewHelper(ts.Logger)

	cfg := &JWT{}

	WithLogger(logger)(cfg)

	assert.Equal(t, logger, cfg.Logger)
}

func TestWithPostFunc(t *testing.T) {
	var (
		cfg          = &JWT{}
		testPostFunc = func(ctx *gin.Context, claims jwt.Claims) error { return nil }
	)

	WithPostFunc(testPostFunc)(cfg)

	assert.NotNil(t, cfg.PostFunc)
}

func TestJWT_Validate(t *testing.T) {
	ts, cleanup, err := tests.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	logger := log.NewHelper(ts.Logger)

	var (
		testKey               = "test key"
		testHeaderName        = "TestHeaderName"
		testHeaderPrefix      = "TestHeaderPrefix"
		testError             = errors.New("test error")
		testPostFuncWithError = func(ctx *gin.Context, claims jwt.Claims) error { return testError }
		testPostFuncNoError   = func(ctx *gin.Context, claims jwt.Claims) error { return nil }
	)

	testInvalidTokenCauseExpired, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)},
		},
	).SignedString([]byte(testKey))
	if err != nil {
		t.Fatal(err)
	}

	testInvalidTokenCauseWrongKey, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Date(3099, 0, 0, 0, 0, 0, 0, time.UTC)},
		},
	).SignedString([]byte("wrong key"))
	if err != nil {
		t.Fatal(err)
	}

	testValidToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Date(2100, 0, 0, 0, 0, 0, 0, time.UTC)},
		},
	).SignedString([]byte(testKey))
	if err != nil {
		t.Fatal(err)
	}

	options := []Option{
		WithLogger(logger),
		WithErrorResponseBody(response.NewBody(int(errorsx.ServerErrorCode), errorsx.ServerErrorCode.String(), nil)),
		WithValidateFailedResponseBody(response.NewBody(int(errorsx.UnauthorizedCode), errorsx.UnauthorizedCode.String(), nil)),
		WithHeaderName(testHeaderName),
		WithHeaderPrefix(testHeaderPrefix),
	}

	type want struct {
		httpStatusCode int
		code           errorsx.ErrorCode
		msg            string
	}

	tts := []struct {
		name     string
		key      string
		token    string
		postFunc func(ctx *gin.Context, claims jwt.Claims) error
		want     want
	}{
		{"without_key", "", "", nil, want{http.StatusInternalServerError, errorsx.ServerErrorCode, ErrMissingKey.Error()}},
		{"without_token", testKey, "", nil, want{http.StatusUnauthorized, errorsx.UnauthorizedCode, errorsx.UnauthorizedCode.String()}},
		{"with_illegal_token", testKey, "eHh4.eXl5.enp6", nil, want{http.StatusUnauthorized, errorsx.UnauthorizedCode, "invalid character 'x' looking for beginning of value"}},
		{"with_expired_token", testKey, testInvalidTokenCauseExpired, nil, want{http.StatusUnauthorized, errorsx.UnauthorizedCode, "Token is expired"}},
		{"with_wrong_token", testKey, testInvalidTokenCauseWrongKey, nil, want{http.StatusUnauthorized, errorsx.UnauthorizedCode, "signature is invalid"}},
		{"call_post_function_error", testKey, testValidToken, testPostFuncWithError, want{http.StatusInternalServerError, errorsx.ServerErrorCode, ErrCallPostFuncFailed.Error()}},
		{"validate_success", testKey, testValidToken, testPostFuncNoError, want{http.StatusOK, errorsx.SuccessCode, errorsx.SuccessCode.String()}},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			ops := append(options, WithPostFunc(tt.postFunc))
			j := New(tt.key, ops...)

			router := setupRouter(j.Validate())

			r := httptest.NewRequest("GET", "/test", nil)
			r.Header.Set(testHeaderName, tt.token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			body := new(response.Body)
			if err := json.Unmarshal(w.Body.Bytes(), body); err != nil {
				t.Fatal(err)
				return
			}

			assert.Equal(t, tt.want.httpStatusCode, w.Code)
			assert.EqualValues(t, tt.want.code, body.Code)
			assert.Equal(t, tt.want.msg, body.Msg)
		})
	}

	t.Run("token_parse_failed", func(t *testing.T) {
		monkey.Patch(jwt.Parse, func(tokenString string, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
			return nil, errors.New("parse failed")
		})
		defer monkey.Unpatch(jwt.Parse)

		j := New(testKey, options...)

		router := setupRouter(j.Validate())

		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set(testHeaderName, "abc")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		body := new(response.Body)
		if err := json.Unmarshal(w.Body.Bytes(), body); err != nil {
			t.Fatal(err)
			return
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.EqualValues(t, errorsx.ServerErrorCode, body.Code)
		assert.Equal(t, ErrParseTokenFailed.Error(), body.Msg)
	})
}

func Test_defaultPostFunc(t *testing.T) {
	ts, cleanup, err := tests.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	logger := log.NewHelper(ts.Logger)

	router := setupRouter(nil)

	var (
		testPath         = "/test-post-func"
		testMethod       = "GET"
		testExceptClaims = jwt.MapClaims{
			"id":   "1",
			"name": "Tom",
		}
	)

	router.Handle(
		testMethod,
		testPath,
		func(ctx *gin.Context) {
			if err := defaultPostFunc(ctx, testExceptClaims); err != nil {
				logger.Error(err)
				return
			}
			ctx.Next()
		},
		func(ctx *gin.Context) {
			claims := ctx.Request.Context().Value(DefaultContextKey).(jwt.MapClaims)

			ctx.JSON(http.StatusOK, gin.H{
				"id":   claims["id"],
				"name": claims["name"],
			})
		},
	)

	r := httptest.NewRequest(testMethod, testPath, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	body := map[string]string{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testExceptClaims["id"], body["id"])
	assert.Equal(t, testExceptClaims["name"], body["name"])
}

func Test_handleResponse(t *testing.T) {
	var (
		testErrorResponseBody          = response.NewBody(int(errorsx.ServerErrorCode), errorsx.ServerErrorCode.String(), nil)
		testValidateFailedResponseBody = response.NewBody(int(errorsx.UnauthorizedCode), errorsx.UnauthorizedCode.String(), nil)
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
				handleResponse(ctx, http.StatusUnauthorized, nil, nil)
				return
			},
			want{http.StatusUnauthorized, 0, ""},
		},
		{
			"with_validate_failed_response_body_without_error",
			func(ctx *gin.Context) {
				handleResponse(ctx, http.StatusUnauthorized, testValidateFailedResponseBody, nil)
				return
			},
			want{http.StatusUnauthorized, errorsx.UnauthorizedCode, errorsx.UnauthorizedCode.String()},
		},
		{
			"with_validate_failed_response_body_with_error",
			func(ctx *gin.Context) {
				handleResponse(ctx, http.StatusUnauthorized, testValidateFailedResponseBody, errors.New("test error"))
				return
			},
			want{http.StatusUnauthorized, errorsx.UnauthorizedCode, "test error"},
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

func setupRouter(mw gin.HandlerFunc) *gin.Engine {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/test", mw, func(c *gin.Context) { response.Success(c); return })

	return r
}
