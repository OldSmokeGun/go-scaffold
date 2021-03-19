package app

import (
	"gin-scaffold/app/utils/validator"
	"github.com/gin-gonic/gin"
)

// FrameInitialize 框架初始化钩子函数
func FrameInitialize() error {
	// 初始化操作

	return nil
}

// ApplicationInitialize 应用初始化钩子函数
func ApplicationInitialize(router *gin.Engine) error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		return err
	}

	return nil
}
