package router

import (
	"go-scaffold/internal/app/transport/http/handler/v1/greet"
	"go-scaffold/internal/app/transport/http/handler/v1/trace"
	"go-scaffold/internal/app/transport/http/handler/v1/user"

	"github.com/gin-gonic/gin"
)

// apiV1Group v1 API routing group
type apiV1Group struct {
	greetHandler greet.HandlerInterface
	traceHandler trace.HandlerInterface
	userHandler  user.HandlerInterface
	group        *gin.RouterGroup
}

// NewAPIV1Group return *apiV1Group
func NewAPIV1Group(
	greetHandler greet.HandlerInterface,
	traceHandler trace.HandlerInterface,
	userHandler user.HandlerInterface,
) *apiV1Group {
	return &apiV1Group{
		greetHandler: greetHandler,
		traceHandler: traceHandler,
		userHandler:  userHandler,
	}
}

func (g *apiV1Group) setup(rg *gin.RouterGroup) {
	g.group = rg.Group("/v1")
}

func (g *apiV1Group) useRoutes() {
	g.group.GET("/greet", g.greetHandler.Hello)
	g.group.GET("/trace", g.traceHandler.Example)

	g.group.GET("/users", g.userHandler.List)
	g.group.GET("/user/:id", g.userHandler.Detail)
	g.group.POST("/user", g.userHandler.Create)
	g.group.PUT("/user/:id", g.userHandler.Update)
	g.group.DELETE("/user/:id", g.userHandler.Delete)
}
