package router

import (
	"gin-scaffold/internal/app/global"
	_ "gin-scaffold/internal/app/rest/api/docs"
	"gin-scaffold/internal/app/rest/config"
	"gin-scaffold/internal/app/rest/handler/greet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
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

	// 覆盖 swagger 配置
	if global.Config().Env != config.Prod {
		// swagger
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	greetHandler := greet.NewHandler()

	router.GET("/greet", greetHandler.Hello)

	// ...

	return router
}
