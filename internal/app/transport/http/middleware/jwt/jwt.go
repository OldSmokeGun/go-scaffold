package jwt

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	pjwt "github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

var (
	ErrMissingKey       = errors.New("key is missing")
	ErrMissingToken     = errors.New("token is missing")
	ErrInvalidToken     = errors.New("token is invalid")
	ErrParseTokenFailed = errors.New("parse token failed")
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
	Debug(...interface{})
	Debugf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

// ResponseBody 响应的 body
type ResponseBody interface {
	WithMsg(msg string)
}

// jwt 中间件的配置项
type jwt struct {

	// key jwt 密钥
	key string

	// headerName jwt HTTP Header 键
	// 如果不指定，默认为 "Authorization"
	headerName string

	// headerPrefix jwt HTTP Header 值前缀
	// 如果不指定，默认为 "Bearer "
	headerPrefix string

	// errorResponseBody 服务器发生错误时以 application/json 方式返回的 body
	// 如果为 nil，则不返回 body
	errorResponseBody ResponseBody

	// validateErrorResponseBody jwt 校验错误时以 application/json 方式返回的 body
	// 如果为 nil，则不返回 body
	validateErrorResponseBody ResponseBody

	// logger 发生错误时记录错误的日志实例
	// 如果为 nil, 则不记录错误
	logger Logger

	// postFunc jwt 校验成功后调用的钩子函数
	// 可以覆盖将 Claims 写入 Context 的默认行为
	postFunc func(ctx *gin.Context, claims pjwt.Claims)

	// raw 原始 token
	raw string
}

// Option jwt 配置选项
type Option func(config *jwt)

// WithHeaderName 设置 jwt HTTP Header 键
func WithHeaderName(headerName string) Option {
	return func(config *jwt) {
		config.headerName = headerName
	}
}

// WithHeaderPrefix 设置 jwt HTTP Header 值前缀
func WithHeaderPrefix(headerPrefix string) Option {
	return func(config *jwt) {
		config.headerPrefix = headerPrefix
	}
}

// WithErrorResponseBody 设置服务器发生错误时以 application/json 方式返回的 body
func WithErrorResponseBody(body ResponseBody) Option {
	return func(config *jwt) {
		config.errorResponseBody = body
	}
}

// WithValidateErrorResponseBody 设置 jwt 校验错误时以 application/json 方式返回的 body
func WithValidateErrorResponseBody(body ResponseBody) Option {
	return func(config *jwt) {
		config.validateErrorResponseBody = body
	}
}

// WithLogger 设置发生错误时记录错误的日志实例
func WithLogger(logger Logger) Option {
	return func(config *jwt) {
		config.logger = logger
	}
}

// WithPostFunc 设置 jwt 校验成功后调用的钩子函数
func WithPostFunc(f func(ctx *gin.Context, claims pjwt.Claims)) Option {
	return func(config *jwt) {
		config.postFunc = f
	}
}

// Validate 验证 jwt
func Validate(key string, options ...Option) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := defaultConfig(key)

		for _, option := range options {
			option(c)
		}

		if c.key == "" {
			errorResponse(ctx, http.StatusInternalServerError, c, ErrMissingKey)
			return
		}

		tokenString := ctx.GetHeader(c.headerName)

		if c.headerPrefix != NoneHeaderPrefix {
			tokenString = strings.TrimPrefix(tokenString, c.headerPrefix)
		}

		if tokenString == "" {
			validateErrorResponse(ctx, http.StatusUnauthorized, c, ErrMissingToken)
			return
		}

		c.raw = tokenString

		if c.logger != nil {
			c.logger.Debugf("token: %s", tokenString)
		}

		token, err := pjwt.Parse(tokenString, func(token *pjwt.Token) (interface{}, error) {
			return []byte(c.key), nil
		})
		if err != nil {
			if c.logger != nil {
				c.logger.Errorf("%s: %s", ErrParseTokenFailed, err)
			}
			errorResponse(ctx, http.StatusInternalServerError, c, ErrParseTokenFailed)
			return
		}

		if !token.Valid {
			validateErrorResponse(ctx, http.StatusUnauthorized, c, ErrInvalidToken)
			return
		}

		if c.postFunc != nil {
			c.postFunc(ctx, token.Claims)
		}

		ctx.Next()
	}
}

func defaultConfig(key string) *jwt {
	return &jwt{
		key:                       key,
		headerName:                defaultHeaderName,
		headerPrefix:              defaultHeaderPrefix,
		errorResponseBody:         nil,
		validateErrorResponseBody: nil,
		logger:                    nil,
		postFunc: func(ctx *gin.Context, claims pjwt.Claims) {
			context.WithValue(ctx.Request.Context(), DefaultContextKey, claims)
		},
		raw: "",
	}
}

// errorResponse 服务器发生错误时的响应
func errorResponse(ctx *gin.Context, code int, c *jwt, err error) {
	if c.errorResponseBody == nil {
		ctx.AbortWithStatus(code)
		return
	}

	body := c.errorResponseBody
	if err != nil {
		body.WithMsg(err.Error())
	}

	ctx.AbortWithStatusJSON(code, body)
	return
}

// validateErrorResponse 校验错误时的响应
func validateErrorResponse(ctx *gin.Context, code int, c *jwt, err error) {
	if c.validateErrorResponseBody == nil {
		ctx.AbortWithStatus(code)
		return
	}

	body := c.validateErrorResponseBody
	if err != nil {
		body.WithMsg(err.Error())
	}

	ctx.AbortWithStatusJSON(code, body)
	return
}
