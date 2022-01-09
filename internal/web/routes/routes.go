package routes

import (
	"fmt"
	"gin-scaffold/internal/web/config"
	"gin-scaffold/internal/web/docs"
	"gin-scaffold/internal/web/global"
	"gin-scaffold/internal/web/handler/greet"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Register 注册路由对象
func Register(router *gin.Engine) {
	var basePath = "/"

	// 覆盖 swagger 配置
	if global.Config().Env != config.Prod {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", global.Config().Host, global.Config().Port)
		docs.SwaggerInfo.BasePath = basePath

		if global.Config().Env == config.Local {
			docs.SwaggerInfo.Schemes = []string{"http"}
		} else if global.Config().Env == config.Test {
			docs.SwaggerInfo.Schemes = []string{"http", "https"}
		}

		// swagger
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	greetHandler := greet.NewHandler()

	router.GET("/greet", greetHandler.Hello)
}
