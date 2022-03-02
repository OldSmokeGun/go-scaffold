package router

import (
	ginzap "github.com/gin-contrib/zap"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/rest/middleware/recover"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"go-scaffold/internal/app/rest/router/api"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"time"

	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// New 返回 gin 路由对象
func New() *gin.Engine {
	output := io.MultiWriter(global.LoggerOutput(), os.Stdout)
	gin.DefaultWriter = output
	gin.DefaultErrorWriter = output
	gin.DisableConsoleColor()

	switch global.Config().Env {
	case "local":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(global.Logger(), time.RFC3339, false))
	router.Use(recover.CustomRecoveryWithZap(global.Logger(), true, func(c *gin.Context, err interface{}) {
		responsex.ServerError(c)
		c.Abort()
	}))

	if global.Config().Trace != nil {
		router.Use(otelgin.Middleware(global.Config().Name))
	}

	// 注册 api 路由组
	api.NewGroup().Registry(router)

	// ...

	return router
}
