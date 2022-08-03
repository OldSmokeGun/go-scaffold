package jwt

import (
	"context"
	"errors"
	"go-scaffold/internal/app/transport/http/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrMissingKey         = errors.New("key is missing")
	ErrParseTokenFailed   = errors.New("parse token failed")
	ErrCallPostFuncFailed = errors.New("call post func failed")
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

// JWT middleware configuration
type JWT struct {

	// Key used for signing
	Key string

	// HeaderName HTTP header key of token
	// if not specified，default: "Authorization"
	HeaderName string

	// HeaderPrefix HTTP header value prefix of token
	// if not specified，default: "Bearer "
	HeaderPrefix string

	// ErrorResponseBody returned in application/json format in case of server error
	// if nil, body is not returned
	ErrorResponseBody middleware.ResponseBody

	// ValidateFailedResponseBody returned in application/json format when validate failed
	// if nil, body is not returned
	ValidateFailedResponseBody middleware.ResponseBody

	// Logger log when an error occurs
	// if nil, no error is logged
	Logger middleware.Logger

	// PostFunc the hook function that will be called after token verification is successful
	// this will override the default behavior of writing Claims to the Context
	PostFunc func(ctx *gin.Context, claims jwt.Claims) error

	// Raw original token
	Raw string
}

// Option middleware JWT function
type Option func(j *JWT)

// WithHeaderName set the HTTP header key of token
func WithHeaderName(headerName string) Option {
	return func(j *JWT) {
		j.HeaderName = headerName
	}
}

// WithHeaderPrefix set the HTTP header value prefix of token
func WithHeaderPrefix(headerPrefix string) Option {
	return func(j *JWT) {
		j.HeaderPrefix = headerPrefix
	}
}

// WithErrorResponseBody body returned in case of server error
func WithErrorResponseBody(body middleware.ResponseBody) Option {
	return func(j *JWT) {
		j.ErrorResponseBody = body
	}
}

// WithValidateFailedResponseBody body returned when validate failed
func WithValidateFailedResponseBody(body middleware.ResponseBody) Option {
	return func(j *JWT) {
		j.ValidateFailedResponseBody = body
	}
}

// WithLogger error Logger
func WithLogger(logger middleware.Logger) Option {
	return func(j *JWT) {
		j.Logger = logger
	}
}

// WithPostFunc set the hook function that will be called after token verification is successful
func WithPostFunc(f func(ctx *gin.Context, claims jwt.Claims) error) Option {
	return func(j *JWT) {
		j.PostFunc = f
	}
}

func New(key string, options ...Option) *JWT {
	j := &JWT{
		Key:                        key,
		HeaderName:                 defaultHeaderName,
		HeaderPrefix:               defaultHeaderPrefix,
		ErrorResponseBody:          nil,
		ValidateFailedResponseBody: nil,
		Logger:                     nil,
		PostFunc:                   defaultPostFunc,
		Raw:                        "",
	}

	for _, opt := range options {
		opt(j)
	}

	return j
}

// Validate check if the token is valid
func (j *JWT) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if j.Key == "" {
			handleResponse(ctx, http.StatusInternalServerError, j.ErrorResponseBody, ErrMissingKey)
			return
		}

		tokenString := ctx.GetHeader(j.HeaderName)

		if j.HeaderPrefix != NoneHeaderPrefix {
			tokenString = strings.TrimPrefix(tokenString, j.HeaderPrefix)
		}

		if tokenString == "" {
			handleResponse(ctx, http.StatusUnauthorized, j.ValidateFailedResponseBody, nil)
			return
		}

		j.Raw = tokenString

		if j.Logger != nil {
			j.Logger.Debugf("token: %s", tokenString)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(j.Key), nil
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				handleResponse(ctx, http.StatusUnauthorized, j.ValidateFailedResponseBody, ve)
				return
			} else {
				if j.Logger != nil {
					j.Logger.Errorf("%s: %s", ErrParseTokenFailed, err)
				}
				handleResponse(ctx, http.StatusInternalServerError, j.ErrorResponseBody, ErrParseTokenFailed)
				return
			}
		}

		if j.PostFunc != nil {
			if err = j.PostFunc(ctx, token.Claims); err != nil {
				if j.Logger != nil {
					j.Logger.Errorf("%s: %s", ErrCallPostFuncFailed, err)
				}
				handleResponse(ctx, http.StatusInternalServerError, j.ErrorResponseBody, ErrCallPostFuncFailed)
				return
			}
		}

		ctx.Next()
	}
}

func defaultPostFunc(ctx *gin.Context, claims jwt.Claims) error {
	claimsContext := context.WithValue(ctx.Request.Context(), DefaultContextKey, claims)
	ctx.Request = ctx.Request.WithContext(claimsContext)
	return nil
}

func handleResponse(ctx *gin.Context, httpStatusCode int, body middleware.ResponseBody, err error) {
	if body == nil {
		ctx.AbortWithStatus(httpStatusCode)
		return
	}

	if err != nil {
		body.WithMsg(err.Error())
	}

	ctx.AbortWithStatusJSON(httpStatusCode, body)
	return
}
