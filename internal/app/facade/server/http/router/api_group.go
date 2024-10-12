package router

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"go-scaffold/internal/app/facade/server/http/api/docs"
	imiddleware "go-scaffold/internal/app/facade/server/http/middleware"
	"go-scaffold/internal/config"
)

// ApiGroup api routing group
type ApiGroup struct {
	env        config.Env
	logger     *slog.Logger
	hsConf     config.HTTPServer
	apiV1Group *ApiV1Group
	group      *echo.Group

	basePath string
}

// NewAPIGroup return *ApiGroup
func NewAPIGroup(
	env config.Env,
	logger *slog.Logger,
	hsConf config.HTTPServer,
	apiV1Group *ApiV1Group,
) *ApiGroup {
	return &ApiGroup{
		env:        env,
		logger:     logger,
		hsConf:     hsConf,
		apiV1Group: apiV1Group,
	}
}

func (g *ApiGroup) setup(prefix string, rg *echo.Group) {
	path := "/api"
	g.group = rg.Group(path)
	g.basePath = prefix + path
}

func (g *ApiGroup) useMiddlewares() {
	// allowed to cross
	g.group.Use(middleware.CORS())
	g.group.Use(imiddleware.Limit(*imiddleware.NewDefaultLimitConfig()))
}

func (g *ApiGroup) useRoutes(e *echo.Echo) {
	// register swagger routing
	g.useSwagger()

	// register v1 version api routing group
	g.apiV1Group.setup(g.basePath, g.group)
	g.apiV1Group.useRoutes()
}

func (g *ApiGroup) useSwagger() {
	// swagger documentation
	if g.env == config.Dev {
		docs.SwaggerInfo.Host = g.hsConf.Addr
		extHost, _ := parseExternalAddr(g.hsConf.ExternalAddr)
		if extHost != "" {
			docs.SwaggerInfo.Host = extHost
		}
		docs.SwaggerInfo.BasePath = g.basePath

		dg := g.group.Group("/docs")
		dg.GET("", func(ctx echo.Context) error {
			return ctx.Redirect(http.StatusMovedPermanently, g.basePath+"/docs/index.html")
		})
		dg.GET("/*", echoSwagger.EchoWrapHandler())
	}
}
