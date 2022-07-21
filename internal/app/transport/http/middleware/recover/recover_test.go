package recover

import (
	"encoding/json"
	"errors"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/tests"
	"go-scaffold/internal/app/transport/http/pkg/response"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCustomRecoveryWithZap(t *testing.T) {
	ts, cleanup, err := tests.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	logger := ts.ZapLogger
	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stderr),
			zapcore.InfoLevel,
		)
	}))

	t.Run("panic_with_broken_pipe_error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := CustomRecoveryWithZap(logger, false, nil)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic_with_broken_pipe", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
		})
	})

	t.Run("panic_with_connection_reset_by_peer_error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := CustomRecoveryWithZap(logger, false, nil)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic_with_connection_reset_by_peer", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
		})
	})

	t.Run("panic_without_stack_without_handler", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := CustomRecoveryWithZap(logger, false, nil)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("panic_with_stack_without_handler", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := CustomRecoveryWithZap(logger, true, nil)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("panic_with_stack_with_handler", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := CustomRecoveryWithZap(logger, true, func(c *gin.Context, err interface{}) {
				response.Error(c, errorsx.ServerError())
				c.Abort()
			})
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			body := new(response.Body)
			if err := json.Unmarshal(w.Body.Bytes(), body); err != nil {
				t.Fatal(err)
				return
			}

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.EqualValues(t, errorsx.ServerErrorCode, body.Code)
			assert.Equal(t, errorsx.ServerErrorCode.String(), body.Msg)
		})
	})
}

func TestRecoveryWithZap(t *testing.T) {
	ts, cleanup, err := tests.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	logger := ts.ZapLogger
	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stderr),
			zapcore.InfoLevel,
		)
	}))

	t.Run("panic_with_broken_pipe_error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := RecoveryWithZap(logger, false)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic_with_broken_pipe", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
		})
	})

	t.Run("panic_with_connection_reset_by_peer_error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := RecoveryWithZap(logger, false)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic_with_connection_reset_by_peer", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
		})
	})

	t.Run("panic_without_stack", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := RecoveryWithZap(logger, false)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})

	t.Run("panic_with_stack", func(t *testing.T) {
		assert.NotPanics(t, func() {
			md := RecoveryWithZap(logger, true)
			router := setupRouter(md)

			r := httptest.NewRequest("GET", "/panic", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})
}

func Test_defaultHandleRecovery(t *testing.T) {
	router := setupRouter(nil)

	var (
		path   = "/test-default-handle-recovery"
		method = "GET"
	)

	router.Handle(
		method,
		path,
		func(ctx *gin.Context) {
			defaultHandleRecovery(ctx, nil)
			ctx.Next()
		},
		func(ctx *gin.Context) {
			panic("test error")
		},
	)

	r := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func setupRouter(mw gin.HandlerFunc) *gin.Engine {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/panic_with_broken_pipe", mw, func(c *gin.Context) {
		panic(&net.OpError{
			Op:  "test request",
			Net: "test network",
			Err: &os.SyscallError{
				Syscall: "test",
				Err:     errors.New("broken pipe"),
			},
		})
	})

	r.GET("/panic_with_connection_reset_by_peer", mw, func(c *gin.Context) {
		panic(&net.OpError{
			Op:  "test request",
			Net: "test network",
			Err: &os.SyscallError{
				Syscall: "test",
				Err:     errors.New("connection reset by peer"),
			},
		})
	})

	r.GET("/panic", mw, func(c *gin.Context) { panic("test error") })

	return r
}
