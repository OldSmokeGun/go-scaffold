package v1

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/rest/api/v1/handler/greet"
	"go-scaffold/internal/app/rest/api/v1/handler/trace"
	"go-scaffold/internal/app/rest/api/v1/handler/user"
)

// Group api v1 路由组
type Group struct {
	BasePath string
}

// NewGroup 构造函数
func NewGroup() *Group {
	return &Group{
		BasePath: "/v1",
	}
}

// Registry 注册路由
func (g Group) Registry(rg *gin.RouterGroup) {
	group := rg.Group(g.BasePath)

	// TODO 编写路由

	var (
		greetHandler = greet.New()
		traceHandler = trace.New()
		userHandler  = user.New()
	)

	group.GET("/greet", greetHandler.Hello)
	group.GET("/trace", traceHandler.Example)

	group.GET("/users", userHandler.List)
	group.GET("/user/:id", userHandler.Detail)
	group.POST("/user", userHandler.Create)
	group.PUT("/user", userHandler.Save)
	group.DELETE("/user/:id", userHandler.Delete)
}
