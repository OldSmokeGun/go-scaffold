package jwt

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	pjwt "github.com/golang-jwt/jwt/v4"
	"go-scaffold/internal/app/transport/http/middleware"
	"net/http"
	"strings"
)

var (
	ErrMissingKey       = errors.New("key is missing")
	ErrInvalidToken     = errors.New("token is invalid")
	ErrParseTokenFailed = errors.New("parse token failed")
)

const (
	// defaultHeaderName default HTTP header key
	defaultHeaderName = "Authorization"

	// defaultHeaderPrefix default HTTP header value prefix
	defaultHeaderPrefix = "Bearer "

	// NoneHeaderPrefix indicates that no HTTP header value prefix
	NoneHeaderPrefix = "-"
)

// DefaultContextKey default Context key
var DefaultContextKey = ContextKey{}

type ContextKey struct{}

// option middleware configuration
type option struct {

	// key used for signing
	key string

	// headerName HTTP header key of token
	// if not specified，default: "Authorization"
	headerName string

	// headerPrefix HTTP header value prefix of token
	// if not specified，default: "Bearer "
	headerPrefix string

	// errorResponseBody returned in application/json format in case of server error
	// if nil, body is not returned
	errorResponseBody middleware.ResponseBody

	// validateFailedResponseBody returned in application/json format when validate failed
	// if nil, body is not returned
	validateFailedResponseBody middleware.ResponseBody

	// logger log when an error occurs
	// if nil, no error is logged
	logger middleware.Logger

	// postFunc the hook function that will be called after token verification is successful
	// this will override the default behavior of writing Claims to the Context
	postFunc func(ctx *gin.Context, claims pjwt.Claims)

	// raw original token
	raw string
}

// Option middleware option function
type Option func(config *option)

// WithHeaderName set the HTTP header key of token
func WithHeaderName(headerName string) Option {
	return func(config *option) {
		config.headerName = headerName
	}
}

// WithHeaderPrefix set the HTTP header value prefix of token
func WithHeaderPrefix(headerPrefix string) Option {
	return func(config *option) {
		config.headerPrefix = headerPrefix
	}
}

// WithErrorResponseBody body returned in case of server error
func WithErrorResponseBody(body middleware.ResponseBody) Option {
	return func(config *option) {
		config.errorResponseBody = body
	}
}

// WithValidateFailedResponseBody body returned when validate failed
func WithValidateFailedResponseBody(body middleware.ResponseBody) Option {
	return func(config *option) {
		config.validateFailedResponseBody = body
	}
}

// WithLogger error logger
func WithLogger(logger middleware.Logger) Option {
	return func(config *option) {
		config.logger = logger
	}
}

// WithPostFunc set the hook function that will be called after token verification is successful
func WithPostFunc(f func(ctx *gin.Context, claims pjwt.Claims)) Option {
	return func(config *option) {
		config.postFunc = f
	}
}

// Validate check if the token is valid
func Validate(key string, options ...Option) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		o := defaultOption(key)

		for _, op := range options {
			op(o)
		}

		if o.key == "" {
			errorResponse(ctx, http.StatusInternalServerError, o, ErrMissingKey)
			return
		}

		tokenString := ctx.GetHeader(o.headerName)

		if o.headerPrefix != NoneHeaderPrefix {
			tokenString = strings.TrimPrefix(tokenString, o.headerPrefix)
		}

		if tokenString == "" {
			validateFailedResponse(ctx, http.StatusUnauthorized, o, nil)
			return
		}

		o.raw = tokenString

		if o.logger != nil {
			o.logger.Debugf("token: %s", tokenString)
		}

		token, err := pjwt.Parse(tokenString, func(token *pjwt.Token) (interface{}, error) {
			return []byte(o.key), nil
		})
		if err != nil {
			if o.logger != nil {
				o.logger.Errorf("%s: %s", ErrParseTokenFailed, err)
			}
			errorResponse(ctx, http.StatusInternalServerError, o, ErrParseTokenFailed)
			return
		}

		if !token.Valid {
			validateFailedResponse(ctx, http.StatusUnauthorized, o, ErrInvalidToken)
			return
		}

		if o.postFunc != nil {
			o.postFunc(ctx, token.Claims)
		}

		ctx.Next()
	}
}

func defaultOption(key string) *option {
	return &option{
		key:                        key,
		headerName:                 defaultHeaderName,
		headerPrefix:               defaultHeaderPrefix,
		errorResponseBody:          nil,
		validateFailedResponseBody: nil,
		logger:                     nil,
		postFunc: func(ctx *gin.Context, claims pjwt.Claims) {
			claimsContext := context.WithValue(ctx.Request.Context(), DefaultContextKey, claims)
			ctx.Request = ctx.Request.WithContext(claimsContext)
		},
		raw: "",
	}
}

func errorResponse(ctx *gin.Context, code int, o *option, err error) {
	if o.errorResponseBody == nil {
		ctx.AbortWithStatus(code)
		return
	}

	body := o.errorResponseBody
	if err != nil {
		body.WithMsg(err.Error())
	}

	ctx.AbortWithStatusJSON(code, body)
	return
}

func validateFailedResponse(ctx *gin.Context, code int, o *option, err error) {
	if o.validateFailedResponseBody == nil {
		ctx.AbortWithStatus(code)
		return
	}

	body := o.validateFailedResponseBody
	if err != nil {
		body.WithMsg(err.Error())
	}

	ctx.AbortWithStatusJSON(code, body)
	return
}
