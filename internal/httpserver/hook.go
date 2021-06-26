package httpserver

import (
	"gin-scaffold/internal/httpserver/appcontext"
	"gin-scaffold/internal/pkg/validator"
	"github.com/gin-gonic/gin"
)

// BeforeRun 是 http 服务启动前的钩子函数
func BeforeRun(router *gin.Engine, appCtx *appcontext.Context) error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		panic(err)
	}

	// ...
	return nil
}
