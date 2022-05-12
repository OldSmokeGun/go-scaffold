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

// option middleware configuration
type option struct {

	// enforcer casbin Enforcer
	enforcer *pcasbin.Enforcer

	requestFunc requestFunc

	// errorResponseBody returned in application/json format in case of server error
	// if nil, body is not returned
	errorResponseBody middleware.ResponseBody

	// validateFailedResponseBody returned in application/json format when validate failed
	// if nil, body is not returned
	validateFailedResponseBody middleware.ResponseBody

	// logger log when an error occurs
	// if nil, no error is logged
	logger middleware.Logger
}

// Option middleware option function
type Option func(config *option)

// WithErrorResponseBody body returned in case of server error
func WithErrorResponseBody(body middleware.ResponseBody) Option {
	return func(o *option) {
		o.errorResponseBody = body
	}
}

// WithValidateFailedResponseBody body returned when validate failed
func WithValidateFailedResponseBody(body middleware.ResponseBody) Option {
	return func(o *option) {
		o.validateFailedResponseBody = body
	}
}

// WithLogger error logger
func WithLogger(logger middleware.Logger) Option {
	return func(o *option) {
		o.logger = logger
	}
}

// Validate check if the request has permission
func Validate(enforcer *pcasbin.Enforcer, rf requestFunc, options ...Option) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		o := &option{}

		for _, op := range options {
			op(o)
		}

		if enforcer == nil {
			errorResponse(ctx, http.StatusInternalServerError, o, ErrNilEnforcer)
			return
		}

		if rf == nil {
			errorResponse(ctx, http.StatusInternalServerError, o, ErrNilRequestFunction)
			return
		}

		o.enforcer = enforcer
		o.requestFunc = rf

		items, err := o.requestFunc(ctx)
		if err != nil {
			if o.logger != nil {
				o.logger.Errorf("%s: %s", ErrGettingCasbinRequestParameters, err)
			}
			errorResponse(ctx, http.StatusInternalServerError, o, ErrGettingCasbinRequestParameters)
			return
		}

		if o.logger != nil {
			o.logger.Debugf("casbin request: %v", items)
		}

		ok, err := o.enforcer.Enforce(items...)
		if err != nil {
			if o.logger != nil {
				o.logger.Errorf("%s: %s", ErrMatchingCasbinRequestParameters, err)
			}
			errorResponse(ctx, http.StatusInternalServerError, o, ErrMatchingCasbinRequestParameters)
			return
		}

		if !ok {
			validateFailedResponse(ctx, http.StatusForbidden, o)
			return
		}

		ctx.Next()
	}
}

func errorResponse(ctx *gin.Context, code int, c *option, err error) {
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

func validateFailedResponse(ctx *gin.Context, code int, c *option) {
	if c.validateFailedResponseBody == nil {
		ctx.AbortWithStatus(code)
		return
	}

	ctx.AbortWithStatusJSON(code, c.validateFailedResponseBody)
	return
}
