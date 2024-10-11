package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	uerr "go-scaffold/pkg/errors"
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

					logger.Error("panic recover", slog.Any("error", err))

					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
