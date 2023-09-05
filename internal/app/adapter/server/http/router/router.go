package router

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	imiddleware "go-scaffold/internal/app/adapter/server/http/middleware"
	"go-scaffold/internal/config"
)

type router struct {
	logger   *slog.Logger
	appName  config.AppName
	appEnv   config.Env
	hsConf   config.HTTPServer
	apiGroup *ApiGroup
}

// New return http router
func New(
	logger *slog.Logger,
	appName config.AppName,
	appEnv config.Env,
	hsConf config.HTTPServer,
	apiGroup *ApiGroup,
) http.Handler {
	r := &router{
		logger:   logger,
		appName:  appName,
		appEnv:   appEnv,
		hsConf:   hsConf,
		apiGroup: apiGroup,
	}

	e := setup(appEnv)
	r.useMiddlewares(e)
	r.useRoutes(e)

	return e
}

func (r *router) useMiddlewares(e *echo.Echo) {
	e.HTTPErrorHandler = imiddleware.ErrorHandler(e.Debug, r.logger)
	e.JSONSerializer = imiddleware.JSONSerializer()
	e.Use(middleware.RequestID())
	e.Use(imiddleware.Recover(r.logger))
	e.Use(imiddleware.Logger(r.logger))
	e.Use(otelecho.Middleware(r.appName.String()))
}

func (r *router) useRoutes(e *echo.Echo) {
	path := ""
	_, extPath := parseExternalAddr(r.hsConf.ExternalAddr)
	if extPath != "" {
		path = "/" + extPath
	}

	group := e.Group(path)
	group.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

	// register api routing group
	r.apiGroup.setup(path, group)
	r.apiGroup.useMiddlewares()
	r.apiGroup.useRoutes(e)
}

func setup(appEnv config.Env) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = appEnv.IsDebug()

	return e
}

// parseExternalAddr parse external address
func parseExternalAddr(addr string) (host, path string) {
	s := strings.SplitN(addr, "/", 2)
	if len(s) == 2 {
		host = s[0]
		path = s[1]
	} else if len(s) == 1 {
		host = s[0]
	}
	return
}
