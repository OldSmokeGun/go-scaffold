package handlers

import (
	"gin-scaffold/internal/ctx"
	"gin-scaffold/internal/logics"
	"gin-scaffold/internal/types"
	"gin-scaffold/internal/utils/http/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(appCtx *ctx.Context) gin.HandlerFunc {
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
