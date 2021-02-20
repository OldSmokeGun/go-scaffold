package app

import (
	"gin-scaffold/app/utils/validator"
	"github.com/gin-gonic/gin"
)

// Initialize APP 启动前的初始化钩子函数
func Initialize(router *gin.Engine) error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		return err
	}

	return nil
}
