package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-scaffold/internal/app/domain"
)

const (
	defaultTokenHeaderKey         = "Authorization"
	defaultTokenHeaderValuePrefix = "Bearer "
)

type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (*domain.UserProfile, error)
}

type TokenRefresher interface {
	RefreshToken(ctx context.Context, userProfile domain.UserProfile, token string) (string, error)
}

type AuthConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// HeaderKey key that get the token from header
	// if not specified，default: "Authorization"
	HeaderKey string

	// HeaderValuePrefix  header value prefix of token
	// if not specified，default: "Bearer "
	HeaderValuePrefix string

	// TokenValidator handle the validate of token
	TokenValidator TokenValidator

	// TokenRefresher handle the refresh of token
	TokenRefresher TokenRefresher
}

func (c *AuthConfig) WithSkipper(skipper middleware.Skipper) *AuthConfig {
	c.Skipper = skipper
	return c
}

func (c *AuthConfig) WithHeaderKey(key string) *AuthConfig {
	c.HeaderKey = key
	return c
}

func (c *AuthConfig) WithHeaderValuePrefix(prefix string) *AuthConfig {
	c.HeaderValuePrefix = prefix
	return c
}

func (c *AuthConfig) WithTokenValidator(handler TokenValidator) *AuthConfig {
	c.TokenValidator = handler
	return c
}

func (c *AuthConfig) WithTokenRefresher(handler TokenRefresher) *AuthConfig {
	c.TokenRefresher = handler
	return c
}

func NewDefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		Skipper:           middleware.DefaultSkipper,
		HeaderKey:         defaultTokenHeaderKey,
		HeaderValuePrefix: defaultTokenHeaderValuePrefix,
	}
}

func Auth(config AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			token := c.Request().Header.Get(config.HeaderKey)

			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			if config.HeaderValuePrefix != "" {
				token = strings.TrimPrefix(token, config.HeaderValuePrefix)
			}

			if config.TokenValidator == nil {
				return next(c)
			}

			user, err := config.TokenValidator.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "malformed token").SetInternal(err)
			}

			if config.TokenRefresher != nil {
				refreshedToken, err := config.TokenRefresher.RefreshToken(c.Request().Context(), *user, token)
				if err != nil {
					return err
				}

				c.Response().Header().Set(config.HeaderKey, refreshedToken)
			}

			return next(&Context{c, *user})
		}
	}
}
