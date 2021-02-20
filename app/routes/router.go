package routes

import (
	"gin-scaffold/app/controllers"
	"github.com/gin-gonic/gin"
)

// Register 函数注册 http 路由
func Register(router *gin.Engine) {
	router.GET("/", controllers.Welcome)
}
