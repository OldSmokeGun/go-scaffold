package web

import (
	"gin-scaffold/internal/pkg/validator"
	"github.com/gin-gonic/gin"
)

// Run 是 http 服务启动前的钩子函数
func Run(router *gin.Engine) error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		panic(err)
	}

	// ...
	return nil
}
