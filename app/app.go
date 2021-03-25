package app

import (
	"fmt"
	"gin-scaffold/app/appcontext"
	appconfig "gin-scaffold/app/config"
	"gin-scaffold/app/routes"
	"gin-scaffold/app/utils/validator"
	"gin-scaffold/internal/components/orm"
	"gin-scaffold/internal/config"
	"gin-scaffold/internal/global"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Initialize 核心初始化钩子函数
func Initialize() error {
	// ...
	return nil
}

// Start 启动 app 的 http 服务
func Start(router *gin.Engine, conf config.Config) {
	var (
		err error
		db  *gorm.DB
	)

	if len(global.GetConfigurator().GetStringMap("db")) > 0 {
		ormConfig := new(orm.Config)
		if err := global.GetConfigurator().UnmarshalKey("db", ormConfig); err != nil {
			panic(err)
		}

		db, err = orm.Initialize(ormConfig)
		if err != nil {
			panic(err)
		}
	}
	// 释放资源
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		if err := sqlDB.Close(); err != nil {
			panic(err)
		}
	}()

	// redis
	// redisConfig := new(redis.Config)
	// if err := global.GetConfigurator().UnmarshalKey("redis", redisConfig); err != nil {
	// 	panic(err)
	// }
	// rc, err := redis.Initialize(redisConfig)
	// if err != nil {
	// 	panic(err)
	// }

	// 创建 app Context 对象
	appCtx := appcontext.NewContext(
		appcontext.WithConfig(appconfig.Config{Config: conf}),
		appcontext.WithDB(db),
		// appcontext.WithRedisClient(rc),
	)

	err = validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		panic(err)
	}

	// 注册路由
	routes.Register(router, appCtx)

	// http 服务启动
	if err := router.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
		panic(err)
	}
}
