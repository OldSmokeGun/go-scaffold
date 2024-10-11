package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type LimitConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// Limiter handle the limit of request
	Limiter *rate.Limiter
}

func (c *LimitConfig) WithSkipper(skipper middleware.Skipper) *LimitConfig {
	c.Skipper = skipper
	return c
}

func (c *LimitConfig) WithLimiter(limiter *rate.Limiter) *LimitConfig {
	c.Limiter = limiter
	return c
}

func NewDefaultLimitConfig() *LimitConfig {
	return &LimitConfig{
		Skipper: middleware.DefaultSkipper,
		Limiter: rate.NewLimiter(rate.Every(time.Second/10), 60),
	}
}

func Limit(config LimitConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if !config.Limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "requests are too frequent")
			}

			return next(c)
		}
	}
}
