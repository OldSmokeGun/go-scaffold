package routes

import (
	"gin-scaffold/app/controllers"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/", controllers.Welcome)
}
