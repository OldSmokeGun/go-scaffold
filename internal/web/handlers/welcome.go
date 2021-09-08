package handlers

import (
	"gin-scaffold/internal/pkg/http/response"
	"gin-scaffold/internal/web/logics"
	"gin-scaffold/internal/web/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Welcome struct{}

func NewWelcome() *Welcome {
	return &Welcome{}
}

func (w *Welcome) Hello(ctx *gin.Context) {
	welcomeLogic := logics.NewWelcomeLogic(ctx)
	req := types.WelcomeReq{}
	resp, err := welcomeLogic.Hello(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ServerError)
	}
	ctx.JSON(http.StatusOK, response.Success.WithData(resp))
}
