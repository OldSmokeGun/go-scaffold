package casbin

import (
	"errors"
	pcasbin "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/transport/http/middleware"
	"net/http"
)

var (
	ErrNilEnforcer                     = errors.New("casbin enforcer is nil pointer")
	ErrNilRequestFunction              = errors.New("casbin request function is nil")
	ErrGettingCasbinRequestParameters  = errors.New("error getting casbin request parameters")
	ErrMatchingCasbinRequestParameters = errors.New("error matching casbin request parameters")
)

// requestFunc how to get request parameters in PERM model
type requestFunc func(ctx *gin.Context) ([]interface{}, error)

// Casbin middleware configuration
type Casbin struct {

	// Enforcer casbin Enforcer
	Enforcer *pcasbin.Enforcer

	RequestFunc requestFunc

	// ErrorResponseBody returned in application/json format in case of server error
	// if nil, body is not returned
	ErrorResponseBody middleware.ResponseBody

	// ValidateFailedResponseBody returned in application/json format when validate failed
	// if nil, body is not returned
	ValidateFailedResponseBody middleware.ResponseBody

	// Logger log when an error occurs
	// if nil, no error is logged
	Logger middleware.Logger
}

// Option middleware Casbin function
type Option func(c *Casbin)

// WithErrorResponseBody body returned in case of server error
func WithErrorResponseBody(body middleware.ResponseBody) Option {
	return func(c *Casbin) {
		c.ErrorResponseBody = body
	}
}

// WithValidateFailedResponseBody body returned when validate failed
func WithValidateFailedResponseBody(body middleware.ResponseBody) Option {
	return func(c *Casbin) {
		c.ValidateFailedResponseBody = body
	}
}

// WithLogger error Logger
func WithLogger(logger middleware.Logger) Option {
	return func(c *Casbin) {
		c.Logger = logger
	}
}

func New(enforcer *pcasbin.Enforcer, rf requestFunc, options ...Option) *Casbin {
	c := &Casbin{
		Enforcer:                   enforcer,
		RequestFunc:                rf,
		ErrorResponseBody:          nil,
		ValidateFailedResponseBody: nil,
		Logger:                     nil,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// Validate check if the request has permission
func (c *Casbin) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if c.Enforcer == nil {
			errorResponse(ctx, http.StatusInternalServerError, c, ErrNilEnforcer)
			return
		}

		if c.RequestFunc == nil {
			errorResponse(ctx, http.StatusInternalServerError, c, ErrNilRequestFunction)
			return
		}

		items, err := c.RequestFunc(ctx)
		if err != nil {
			if c.Logger != nil {
				c.Logger.Errorf("%s: %s", ErrGettingCasbinRequestParameters, err)
			}
			errorResponse(ctx, http.StatusInternalServerError, c, ErrGettingCasbinRequestParameters)
			return
		}

		if c.Logger != nil {
			c.Logger.Debugf("casbin request: %v", items)
		}

		ok, err := c.Enforcer.Enforce(items...)
		if err != nil {
			if c.Logger != nil {
				c.Logger.Errorf("%s: %s", ErrMatchingCasbinRequestParameters, err)
			}
			errorResponse(ctx, http.StatusInternalServerError, c, ErrMatchingCasbinRequestParameters)
			return
		}

		if !ok {
			validateFailedResponse(ctx, http.StatusForbidden, c)
			return
		}

		ctx.Next()
	}
}

func errorResponse(ctx *gin.Context, code int, c *Casbin, err error) {
	if c.ErrorResponseBody == nil {
		ctx.AbortWithStatus(code)
		return
	}

	body := c.ErrorResponseBody
	if err != nil {
		body.WithMsg(err.Error())
	}

	ctx.AbortWithStatusJSON(code, body)
	return
}

func validateFailedResponse(ctx *gin.Context, code int, c *Casbin) {
	if c.ValidateFailedResponseBody == nil {
		ctx.AbortWithStatus(code)
		return
	}

	ctx.AbortWithStatusJSON(code, c.ValidateFailedResponseBody)
	return
}
