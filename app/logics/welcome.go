package logics

import (
	"gin-scaffold/app/appcontext"
	"gin-scaffold/app/types"
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
