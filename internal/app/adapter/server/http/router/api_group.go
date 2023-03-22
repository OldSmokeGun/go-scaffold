package router

import (
	"net/http"

	"go-scaffold/internal/app/adapter/server/http/api/docs"
	"go-scaffold/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/exp/slog"
)

// ApiGroup api routing group
type ApiGroup struct {
	env      config.Env
	logger   *slog.Logger
	httpConf config.HTTP
	// jwtConf    config.JWT       // optional jwt middleware
	// enforcer   *casbin.Enforcer // optional casbin middleware
	apiV1Group *ApiV1Group
	group      *echo.Group

	basePath string
}

// NewAPIGroup return *ApiGroup
func NewAPIGroup(
	env config.Env,
	logger *slog.Logger,
	httpConf config.HTTP,
	// jwtConf config.JWT, // optional jwt middleware
	// enforcer *casbin.Enforcer, // optional casbin middleware
	apiV1Group *ApiV1Group,
) *ApiGroup {
	return &ApiGroup{
		env:      env,
		logger:   logger,
		httpConf: httpConf,
		// jwtConf:    jwtConf,  // optional jwt middleware
		// enforcer:   enforcer, // optional casbin middleware
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

	// g.group.Use(jwtmw.JWT([]byte(g.jwtConf.Key))) // optional jwt middleware
	// optional casbin middleware
	// if err := g.enforcer.LoadPolicy(); err != nil {
	// 	g.logger.Error("load casbin policy error", err)
	// } else {
	// 	g.group.Use(casbinmw.Middleware(g.enforcer))
	// }
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
		docs.SwaggerInfo.Host = g.httpConf.Addr
		extHost, _ := parseExternalAddr(g.httpConf.ExternalAddr)
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
