package app

import (
	"gin-scaffold/app/utils/validator"
	"github.com/gin-gonic/gin"
)

// app 启动前的钩子函数
func Run(r *gin.Engine) error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		return err
	}

	return nil
}
