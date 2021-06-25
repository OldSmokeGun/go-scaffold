package internal

import (
	"fmt"
	"gin-scaffold/internal/ctx"
	"gin-scaffold/internal/routes"
	"gin-scaffold/internal/utils/validator"
	"gin-scaffold/pkg/orm"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func Initialize(appCtx *ctx.Context) {
	var (
		err error
		db  *gorm.DB
	)

	if len(appCtx.GetConfigurator().GetStringMap("DB")) > 0 {
		ormConfig := new(orm.Config)
		if err := appCtx.GetConfigurator().UnmarshalKey("DB", ormConfig); err != nil {
			panic(err)
		}

		db, err = orm.Initialize(ormConfig)
		if err != nil {
			panic(err)
		}
	}

	// redis
	// redisConfig := new(redis.Config)
	// if err := appCtx.GetConfigurator().UnmarshalKey("redis", redisConfig); err != nil {
	// 	panic(err)
	// }
	// rc, err := redis.Initialize(redisConfig)
	// if err != nil {
	// 	panic(err)
	// }

	appCtx.SetDB(db)
	// appCtx.SetRedisClient(rc)
	// ...
}

// Start 启动 app 的 http 服务
func Start(cmd *cobra.Command, appCtx *ctx.Context) {
	// 释放数据库等资源
	defer func(appCtx *ctx.Context) {
		sqlDB, err := appCtx.GetDB().DB()
		if err != nil {
			panic(err)
		}

		if err := sqlDB.Close(); err != nil {
			panic(err)
		}

		if err := appCtx.GetRedisClient().Close(); err != nil {
			panic(err)
		}
	}(appCtx)

	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		panic(err)
	}

	// 注册路由
	router := gin.Default()
	routes.Register(router, appCtx)

	// http 服务启动
	if err := router.Run(fmt.Sprintf("%s:%d", appCtx.Config.AppConf.Host, appCtx.Config.AppConf.Port)); err != nil {
		panic(err)
	}
}
