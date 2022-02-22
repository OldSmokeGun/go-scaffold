package router

import (
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/rest/router/api"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// New 返回 gin 路由对象
func New() *gin.Engine {
	gin.DefaultWriter = io.MultiWriter(global.LoggerOutput(), os.Stdout)

	switch global.Config().Env {
	case "local":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	if global.Config().Trace != nil {
		router.Use(otelgin.Middleware(global.Config().Name))
	}

	// 注册 api 路由组
	api.NewGroup().Registry(router)

	// ...

	return router
}
