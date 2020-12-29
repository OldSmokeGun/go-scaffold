package router

import (
	"gin-scaffold/app/controller"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/", controller.Welcome)
}
