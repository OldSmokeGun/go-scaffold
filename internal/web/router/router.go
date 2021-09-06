package router

import (
	"gin-scaffold/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

// Register 注册路由对象
func Register(handler *gin.Engine) {
	Welcome := handlers.NewWelcome()

	handler.GET("/", Welcome.Hello)
}
