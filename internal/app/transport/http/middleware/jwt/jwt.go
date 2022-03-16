package jwt

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

var (
	ErrNotProvideKey    = errors.New("未提供 key")
	ErrFailedToGetToken = errors.New("未能获取到 token")
	ErrInvalidToken     = errors.New("无效的 token")
)

const (
	// defaultHeaderName 默认的 HTTP Header 键
	defaultHeaderName = "Authorization"

	// defaultHeaderPrefix 默认的 HTTP Header 值前缀
	defaultHeaderPrefix = "Bearer "

	// NoneHeaderPrefix 表示不指定 HTTP Header 值前缀
	NoneHeaderPrefix = "-"
)

// DefaultContextKey 默认的 Context 键
var DefaultContextKey = ContextKey{}

type ContextKey struct{}

// Logger 错误日志实例需要实现的接口
type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
}

// JWT 中间件的配置项
type JWT struct {

	// Key JWT 密钥
	Key string

	// HeaderName JWT HTTP Header 键
	// 如果不指定，默认为 "Authorization"
	HeaderName string

	// HeaderPrefix JWT HTTP Header 值前缀
	// 如果不指定，默认为 "Bearer "
	HeaderPrefix string

	// ErrorResponseBody 服务器发生错误时以 application/json 方式返回的 body
	// 如果为 nil，则不返回 body
	ErrorResponseBody interface{}

	// ValidateErrorResponseBody JWT 校验错误时以 application/json 方式返回的 body
	// 如果为 nil，则不返回 body
	ValidateErrorResponseBody interface{}

	// Logger 发生错误时记录错误的日志实例
	// 如果为 nil, 则不记录错误
	Logger Logger

	// PostFunc JWT 校验成功后调用的钩子函数
	// 可以覆盖将 Claims 写入 Context 的默认行为
	PostFunc func(ctx *gin.Context, claims jwt.Claims)

	// raw 原始 token
	raw string
}

// Option JWT 配置选项
type Option func(config *JWT)

// WithHeaderName 设置 JWT HTTP Header 键
func WithHeaderName(headerName string) Option {
	return func(config *JWT) {
		config.HeaderName = headerName
	}
}

// WithHeaderPrefix 设置 JWT HTTP Header 值前缀
func WithHeaderPrefix(headerPrefix string) Option {
	return func(config *JWT) {
		config.HeaderPrefix = headerPrefix
	}
}

// WithErrorResponseBody 设置服务器发生错误时以 application/json 方式返回的 body
func WithErrorResponseBody(body interface{}) Option {
	return func(config *JWT) {
		config.ErrorResponseBody = body
	}
}

// WithValidateErrorResponseBody 设置 JWT 校验错误时以 application/json 方式返回的 body
func WithValidateErrorResponseBody(body interface{}) Option {
	return func(config *JWT) {
		config.ValidateErrorResponseBody = body
	}
}

// WithLogger 设置发生错误时记录错误的日志实例
func WithLogger(logger Logger) Option {
	return func(config *JWT) {
		config.Logger = logger
	}
}

// WithPostFunc 设置 JWT 校验成功后调用的钩子函数
func WithPostFunc(f func(ctx *gin.Context, claims jwt.Claims)) Option {
	return func(config *JWT) {
		config.PostFunc = f
	}
}

// Auth 验证 JWT
func Auth(key string, options ...Option) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := defaultConfig(key)

		for _, option := range options {
			option(c)
		}

		if c.Key == "" {
			errorHandle(ctx, c, ErrNotProvideKey)
			return
		}

		tokenString := ctx.GetHeader(c.HeaderName)

		if c.HeaderPrefix != NoneHeaderPrefix {
			tokenString = strings.TrimPrefix(tokenString, c.HeaderPrefix)
		}

		if tokenString == "" {
			validateErrorHandle(ctx, c, ErrFailedToGetToken)
			return
		}

		c.raw = tokenString

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(c.Key), nil
		})
		if err != nil {
			validateErrorHandle(ctx, c, err)
			return
		}

		if !token.Valid {
			validateErrorHandle(ctx, c, ErrInvalidToken)
			return
		}

		if c.PostFunc != nil {
			c.PostFunc(ctx, token.Claims)
		}
	}
}

func defaultConfig(key string) *JWT {
	return &JWT{
		Key:                       key,
		HeaderName:                defaultHeaderName,
		HeaderPrefix:              defaultHeaderPrefix,
		ErrorResponseBody:         nil,
		ValidateErrorResponseBody: nil,
		Logger:                    nil,
		PostFunc: func(ctx *gin.Context, claims jwt.Claims) {
			context.WithValue(ctx.Request.Context(), DefaultContextKey, claims)
		},
		raw: "",
	}
}

// errorHandle 服务器发生错误时的操作
func errorHandle(ctx *gin.Context, c *JWT, err error) {
	if c.Logger != nil {
		c.Logger.Error(err)
	}

	if c.ErrorResponseBody == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, c.ErrorResponseBody)
	return
}

// validateErrorHandle 校验错误时的操作
func validateErrorHandle(ctx *gin.Context, c *JWT, err error) {
	if c.Logger != nil {
		c.Logger.Errorf("token: %s -> %s", c.raw, err.Error())
	}

	if c.ValidateErrorResponseBody == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, c.ValidateErrorResponseBody)
	return
}
