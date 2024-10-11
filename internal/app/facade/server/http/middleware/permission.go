package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type PermissionValidator interface {
	ValidatePermission(ctx context.Context, user int64, permissionKey string) (bool, error)
}

type PermissionConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// PermissionValidator handle the validate of permission
	PermissionValidator PermissionValidator
}

func (c *PermissionConfig) WithSkipper(skipper middleware.Skipper) *PermissionConfig {
	c.Skipper = skipper
	return c
}

func (c *PermissionConfig) WithValidator(handler PermissionValidator) *PermissionConfig {
	c.PermissionValidator = handler
	return c
}

func NewDefaultPermissionConfig() *PermissionConfig {
	return &PermissionConfig{
		Skipper: middleware.DefaultSkipper,
	}
}

func Permission(config PermissionConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if config.PermissionValidator != nil {
				user := c.(*Context).GetUser()
				permissionKey := fmt.Sprintf("%s %s", c.Request().Method, c.Path())

				result, err := config.PermissionValidator.ValidatePermission(c.Request().Context(), user.ID, permissionKey)
				if err != nil {
					return err
				}
				if !result {
					return echo.NewHTTPError(http.StatusForbidden, "access denied")
				}
			}

			return next(c)
		}
	}
}
