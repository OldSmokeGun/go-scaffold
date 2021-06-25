package logics

import (
	"gin-scaffold/internal/ctx"
	"gin-scaffold/internal/types"
	"github.com/gin-gonic/gin"
)

type WelcomeLogic struct {
	ctx    *gin.Context
	appCtx *ctx.Context
}

func NewWelcomeLogic(ctx *gin.Context, appCtx *ctx.Context) *WelcomeLogic {
	return &WelcomeLogic{
		ctx:    ctx,
		appCtx: appCtx,
	}
}

func (l WelcomeLogic) Welcome(req types.WelcomeReq) (types.WelcomeResp, error) {
	return types.WelcomeResp{
		Msg: "Welcome",
	}, nil
}
