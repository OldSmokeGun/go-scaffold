package router

import (
	"fmt"
	"gin-scaffold/internal/app/config"
	"gin-scaffold/internal/app/global"
	"gin-scaffold/internal/app/rest/api/docs"

	"gin-scaffold/internal/app/rest/handler/greet"
	"gin-scaffold/internal/app/rest/pkg/swagger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"os"
)

// New 返回 gin 路由对象
func New() *gin.Engine {
	gin.DefaultWriter = io.MultiWriter(global.LoggerOutput(), os.Stdout)

	switch global.Config().Env {
	case config.Local:
		gin.SetMode(gin.DebugMode)
	case config.Test:
		gin.SetMode(gin.TestMode)
	case config.Prod:
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 允许跨越
	router.Use(cors.Default())

	// swagger 配置
	if global.Config().Env != config.Prod {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", docs.SwaggerInfo.Host, global.Config().REST.Port)
		if global.Config().REST.ExternalUrl != "" {
			docs.SwaggerInfo.Host = global.Config().REST.ExternalUrl
		}

		swagger.Setup(router, swagger.Config{Path: "/docs", OptionFunc: func(c *ginSwagger.Config) {
			c.DefaultModelsExpandDepth = -1
		}})
	}

	greetHandler := greet.NewHandler()

	router.GET("/greet", greetHandler.Hello)

	// ...

	return router
}
