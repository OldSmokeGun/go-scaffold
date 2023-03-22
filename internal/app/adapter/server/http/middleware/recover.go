package middleware

import (
	"fmt"
	"net/http"

	uerr "go-scaffold/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

// Recover returns a recover middleware
func Recover(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					if r == http.ErrAbortHandler {
						panic(r)
					}

					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					if !uerr.IsStackTrace(err) {
						err = errors.WithStack(err)
					}

					logger.Error("panic recover", err)

					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
