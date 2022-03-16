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

// DefaultContextKey 默认的 Context 键
var DefaultContextKey = ContextKey{}

const (
	// DefaultHeaderName 默认的 HTTP Header 键
	DefaultHeaderName = "Authorization"

	// DefaultHeaderPrefix 默认的 HTTP Header 值前缀
	DefaultHeaderPrefix = "Bearer "

	// NoneHeaderPrefix 表示不指定 HTTP Header 值前缀
	NoneHeaderPrefix = "-"
)

type ContextKey struct{}

// Logger 错误日志实例需要实现的接口
type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
}

// Config 中间件的配置项
type Config struct {

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

	// ContextKey JWT 校验成功后，写入到 Context 的 key
	ContextKey interface{}

	// raw 原始 token
	raw string
}

// Option JWT 配置选项
type Option func(config *Config)

// WithHeaderName 设置 JWT HTTP Header 键
func WithHeaderName(headerName string) Option {
	return func(config *Config) {
		config.HeaderName = headerName
	}
}

// WithHeaderPrefix 设置 JWT HTTP Header 值前缀
func WithHeaderPrefix(headerPrefix string) Option {
	return func(config *Config) {
		config.HeaderPrefix = headerPrefix
	}
}

// WithErrorResponseBody 设置服务器发生错误时以 application/json 方式返回的 body
func WithErrorResponseBody(body interface{}) Option {
	return func(config *Config) {
		config.ErrorResponseBody = body
	}
}

// WithValidateErrorResponseBody 设置 JWT 校验错误时以 application/json 方式返回的 body
func WithValidateErrorResponseBody(body interface{}) Option {
	return func(config *Config) {
		config.ValidateErrorResponseBody = body
	}
}

// WithLogger 设置发生错误时记录错误的日志实例
func WithLogger(logger Logger) Option {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithContextKey 设置 JWT 校验成功后，写入到 Context 的 key
func WithContextKey(key interface{}) Option {
	return func(config *Config) {
		config.ContextKey = key
	}
}

// Auth 验证 JWT
func Auth(key string, options ...Option) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config := defaultConfig(key)

		for _, option := range options {
			option(config)
		}

		if config.Key == "" {
			errorHandle(ctx, config, ErrNotProvideKey)
			return
		}

		tokenString := ctx.GetHeader(config.HeaderName)

		if config.HeaderPrefix != NoneHeaderPrefix {
			tokenString = strings.TrimPrefix(tokenString, config.HeaderPrefix)
		}

		if tokenString == "" {
			validateErrorHandle(ctx, config, ErrFailedToGetToken)
			return
		}

		config.raw = tokenString

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Key), nil
		})
		if err != nil {
			validateErrorHandle(ctx, config, err)
			return
		}

		if !token.Valid {
			validateErrorHandle(ctx, config, ErrInvalidToken)
			return
		}

		context.WithValue(ctx.Request.Context(), config.ContextKey, token.Claims)
	}
}

func defaultConfig(key string) *Config {
	return &Config{
		Key:                       key,
		HeaderName:                DefaultHeaderName,
		HeaderPrefix:              DefaultHeaderPrefix,
		ErrorResponseBody:         nil,
		ValidateErrorResponseBody: nil,
		Logger:                    nil,
		ContextKey:                DefaultContextKey,
		raw:                       "",
	}
}

// errorHandle 服务器发生错误时的操作
func errorHandle(ctx *gin.Context, c *Config, err error) {
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
func validateErrorHandle(ctx *gin.Context, c *Config, err error) {
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
