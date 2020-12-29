package controller

import (
	"github.com/gin-gonic/gin"
)

func Welcome(c *gin.Context) {
	c.String(200, "<h1>Welcome</h1>")
}
