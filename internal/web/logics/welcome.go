package logics

import (
	"gin-scaffold/internal/web/types"
	"github.com/gin-gonic/gin"
)

type WelcomeLogic struct {
	ctx *gin.Context
}

func NewWelcomeLogic(ctx *gin.Context) *WelcomeLogic {
	return &WelcomeLogic{
		ctx: ctx,
	}
}

func (l *WelcomeLogic) Hello(req types.WelcomeReq) (types.WelcomeResp, error) {
	return types.WelcomeResp{
		Msg: "Hello",
	}, nil
}
