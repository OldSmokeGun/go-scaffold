package router

import (
	"gin-scaffold/global"
	"gin-scaffold/internal/web/handlers"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

// Setup 函数创建 http 路由对象
func Setup() *gin.Engine {
	gin.DefaultWriter = io.MultiWriter(global.LogRotate(), os.Stdout)

	router := gin.Default()

	Welcome := handlers.NewWelcome()

	router.GET("/", Welcome.Hello)

	return router
}
