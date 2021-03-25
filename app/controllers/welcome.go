package controllers

import (
	"gin-scaffold/app/appcontext"
	"github.com/gin-gonic/gin"
)

func Welcome(appCtx *appcontext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(200, "<h1>Welcome</h1>")
	}
}
