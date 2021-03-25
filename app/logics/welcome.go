package logics

import (
	"gin-scaffold/app/appcontext"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WelcomeLogic struct {
	ctx    *gin.Context
	appCtx *appcontext.Context
}

func (l WelcomeLogic) Welcome() {
	l.ctx.String(http.StatusOK, "<h1>Welcome</h1>")
}
