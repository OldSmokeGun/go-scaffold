package router

import (
	"gin-scaffold/internal/httpserver/appcontext"
	"gin-scaffold/internal/httpserver/handlers"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// Setup 函数创建 http 路由对象
func Setup(appCtx *appcontext.Context) *gin.Engine {
	gin.DefaultWriter = io.MultiWriter(appCtx.LogRotate(), os.Stdout)

	router := gin.Default()

	router.GET("/", handlers.Welcome(appCtx))

	return router
}
