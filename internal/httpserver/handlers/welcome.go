package handlers

import (
	"gin-scaffold/internal/httpserver/appcontext"
	"gin-scaffold/internal/httpserver/logics"
	"gin-scaffold/internal/httpserver/types"
	"gin-scaffold/internal/pkg/http/response"
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
