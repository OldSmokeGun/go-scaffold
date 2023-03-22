package router

import (
	v1 "go-scaffold/internal/app/adapter/server/http/handler/v1"

	"github.com/labstack/echo/v4"
)

// ApiV1Group v1 API routing group
type ApiV1Group struct {
	greetHandler *v1.GreetHandler
	traceHandler *v1.TraceHandler
	userHandler  *v1.UserHandler
	group        *echo.Group

	basePath string
}

// NewAPIV1Group return *ApiV1Group
func NewAPIV1Group(
	greetHandler *v1.GreetHandler,
	traceHandler *v1.TraceHandler,
	userHandler *v1.UserHandler,
) *ApiV1Group {
	return &ApiV1Group{
		greetHandler: greetHandler,
		traceHandler: traceHandler,
		userHandler:  userHandler,
	}
}

func (g *ApiV1Group) setup(prefix string, rg *echo.Group) {
	path := "/v1"
	g.group = rg.Group(path)
	g.basePath = prefix + path
}

func (g *ApiV1Group) useRoutes() {
	g.group.GET("/greet", g.greetHandler.Hello)
	g.group.GET("/trace", g.traceHandler.Example)

	g.group.GET("/users", g.userHandler.List)
	g.group.GET("/user/:id", g.userHandler.Detail)
	g.group.POST("/user", g.userHandler.Create)
	g.group.PUT("/user/:id", g.userHandler.Update)
	g.group.DELETE("/user/:id", g.userHandler.Delete)
}
