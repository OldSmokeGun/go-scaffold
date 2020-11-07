package router

import (
	"gin-scaffold/internal/app/controllers"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/", controllers.Index)
}
