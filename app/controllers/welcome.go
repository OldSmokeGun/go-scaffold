package controllers

import (
	"gin-scaffold/app/appcontext"
	"gin-scaffold/app/logics"
	"gin-scaffold/app/types"
	"gin-scaffold/app/utils/http/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(appCtx *appcontext.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		welcomeLogic := logics.NewWelcomeLogic(ctx, appCtx)
		req := types.WelcomeReq{}
		resp, err := welcomeLogic.Welcome(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		ctx.JSON(http.StatusOK, response.SuccessFormat(resp))
	}
}
