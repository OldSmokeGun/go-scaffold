package routers

import (
	"gin-scaffold/kernel/app/controller"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/", controller.Index)
}
