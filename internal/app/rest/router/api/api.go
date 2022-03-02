package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/rest/handler/docs"
	"go-scaffold/internal/app/rest/pkg/swagger"
	v1 "go-scaffold/internal/app/rest/router/api/v1"
)

// Group api 路由组
type Group struct {
	BasePath string
	Config   *config.Config
}

// NewGroup 构造函数
func NewGroup() *Group {
	return &Group{
		BasePath: "/api",
		Config:   global.Config(),
	}
}

// Registry 注册路由
func (g Group) Registry(router *gin.Engine) {
	group := router.Group(g.BasePath)
	group.Use(
		cors.Default(), // 允许跨越
		// jwt.Auth(jwt.Config{
		// 	Key:                       global.Config().REST.Jwt.Key,
		// 	ErrorResponseBody:         responsex.NewServerErrorBody(),
		// 	ValidateErrorResponseBody: responsex.NewUnauthorizedBody(),
		// 	Logger:                    global.Logger().Sugar(),
		// 	ContextKey:                "AuthInfo",
		// }), // jwt 认证
	)
	{
		// swagger 配置
		if g.Config.Env != "prod" {
			docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", docs.SwaggerInfo.Host, g.Config.REST.Port)
			if g.Config.REST.ExternalUrl != "" {
				docs.SwaggerInfo.Host = g.Config.REST.ExternalUrl
			}
			docs.SwaggerInfo.BasePath = group.BasePath()

			swagger.Setup(router, swagger.Config{
				Path: group.BasePath() + "/docs",
				OptionFunc: func(c *ginSwagger.Config) {
					c.DefaultModelsExpandDepth = -1
				},
			})
		}

		// 注册 v1 版本路由组
		v1.NewGroup().Registry(group)
	}
}
