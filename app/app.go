package app

import (
	"gin-scaffold/app/utils/validator"
	"github.com/gin-gonic/gin"
)

// CoreInitialize 核心初始化钩子函数
func CoreInitialize() error {
	// 初始化操作
	// redisConfig := new(redis.Config)
	// if err := global.Configurator().UnmarshalKey("redis", redisConfig); err != nil {
	// 	panic(err)
	// }
	// rc, err := redis.Initialize(redisConfig)
	// if err != nil {
	// 	return err
	// }
	// SetRedisClient(rc)

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
