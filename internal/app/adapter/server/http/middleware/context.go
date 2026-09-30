package middleware

import (
	"github.com/labstack/echo/v4"

	"go-scaffold/internal/app/domain"
)

// Context user profile context
type Context struct {
	echo.Context
	user domain.SystemUserProfile
}

// GetUser get user profile from context
func (u *Context) GetUser() domain.SystemUserProfile {
	return u.user
}

// SetUser set user profile to context
func (u *Context) SetUser(user domain.SystemUserProfile) {
	u.user = user
}
