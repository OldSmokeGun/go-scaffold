package jwt

import (
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
	// DefaultHeaderName 默认的 HTTP Header 键
	DefaultHeaderName = "Authorization"

	// DefaultHeaderPrefix 默认的 HTTP Header 值前缀
	DefaultHeaderPrefix = "Bearer "

	// NoneHeaderPrefix 表示不指定 HTTP Header 值前缀
	NoneHeaderPrefix = "-"
)

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

	// ContextFiled 指定 JWT 校验成功后，写入到 gin 上下文的 key
	ContextKey string

	// raw 原始 token
	raw string
}

// Auth 验证 JWT
func Auth(c Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if c.Key == "" {
			errorHandle(ctx, c, ErrNotProvideKey)
			return
		}

		if c.HeaderName == "" {
			c.HeaderName = DefaultHeaderName
		}

		tokenString := ctx.GetHeader(c.HeaderName)

		if c.HeaderPrefix != NoneHeaderPrefix {
			if c.HeaderPrefix == "" {
				tokenString = strings.TrimPrefix(tokenString, DefaultHeaderPrefix)
			} else {
				tokenString = strings.TrimPrefix(tokenString, c.HeaderPrefix)
			}
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

		if c.ContextKey != "" {
			ctx.Set(c.ContextKey, token.Claims)
		}
	}
}

// errorHandle 服务器发生错误时的操作
func errorHandle(ctx *gin.Context, c Config, err error) {
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
func validateErrorHandle(ctx *gin.Context, c Config, err error) {
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
