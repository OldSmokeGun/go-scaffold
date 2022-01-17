package v1

import (
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/rest/api/v1/handler/greet"
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

	greetHandler := greet.NewHandler()

	group.GET("/greet", greetHandler.Hello)
}
