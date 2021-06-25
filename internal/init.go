package internal

import (
	"fmt"
	appcontext "gin-scaffold/internal/context"
	"gin-scaffold/internal/routes"
	"gin-scaffold/internal/utils/validator"
	"gin-scaffold/pkg/orm"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"strconv"
)

func initialize(cmd *cobra.Command, appCtx *appcontext.Context) {
	var (
		err  error
		host string
		port int
		flag = cmd.Flags()
	)

	hostFlag := flag.Lookup("host")
	portFlag := flag.Lookup("port")

	if hostFlag.Changed {
		host = hostFlag.Value.String()
	} else {
		if appCtx.GetConfigurator().InConfig("Host") {
			host = appCtx.GetConfigurator().GetString("Host")
		} else {
			host = hostFlag.DefValue
		}
	}

	if portFlag.Changed {
		port, err = strconv.Atoi(portFlag.Value.String())
		if err != nil {
			panic(err)
		}
	} else {
		if appCtx.GetConfigurator().InConfig("Port") {
			port = appCtx.GetConfigurator().GetInt("Port")
		} else {
			port, err = strconv.Atoi(portFlag.DefValue)
			if err != nil {
				panic(err)
			}
		}
	}

	// 设置 conf 对象中的属性
	appCtx.Config.AppConf.Host = host
	appCtx.Config.AppConf.Port = port
}

// Start 启动 app 的 http 服务
func Start(cmd *cobra.Command, appCtx *appcontext.Context) {
	initialize(cmd, appCtx)

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
	// if err := appCtx.GetConfigurator().UnmarshalKey("redis", redisConfig); err != nil {
	// 	panic(err)
	// }
	// rc, err := redis.Initialize(redisConfig)
	// if err != nil {
	// 	panic(err)
	// }

	appCtx.SetDB(db)

	err = validator.RegisterValidator([]validator.CustomValidator{
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
