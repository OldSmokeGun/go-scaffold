package logics

import (
	appcontext "gin-scaffold/internal/context"
	"gin-scaffold/internal/types"
	"github.com/gin-gonic/gin"
)

type WelcomeLogic struct {
	ctx    *gin.Context
	appCtx *appcontext.Context
}

func NewWelcomeLogic(ctx *gin.Context, appCtx *appcontext.Context) *WelcomeLogic {
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
