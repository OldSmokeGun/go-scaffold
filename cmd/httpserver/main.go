package main

import (
	"fmt"
	"gin-scaffold/global"
	"gin-scaffold/internal/httpserver"
	"gin-scaffold/internal/httpserver/appconfig"
	"gin-scaffold/internal/httpserver/appcontext"
	"gin-scaffold/internal/httpserver/routes"
	"gin-scaffold/pkg/configurator"
	"gin-scaffold/pkg/logger"
	"gin-scaffold/pkg/orm"
	"gin-scaffold/pkg/redisclient"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
	"path/filepath"
)

// defaultConfigPath 默认配置文件路径
var defaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "config/httpserver.yaml")

func main() {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", defaultConfigPath, "配置文件路径")
	pflag.Parse()

	// 加载配置
	var conf appconfig.Config
	if err := configurator.LoadConfig(configPath, &conf); err != nil {
		panic(err)
	}

	var err error

	// 日志初始化
	var log *logrus.Logger
	if conf.LogConf != nil {
		log, err = logger.Setup(*conf.LogConf)
		if err != nil {
			panic(err)
		}
	}

	// orm 初始化
	var db *gorm.DB
	if conf.DatabaseConf != nil {
		db, err = orm.Setup(*conf.DatabaseConf)
		if err != nil {
			panic(err)
		}
	}

	// redis 初始化
	var rc *redis.Client
	if conf.RedisConf != nil {
		rc, err = redisclient.Setup(*conf.RedisConf)
		if err != nil {
			panic(err)
		}
	}

	// 创建上下文依赖
	appCtx := appcontext.New()
	appCtx.SetConfig(conf)
	appCtx.SetLogger(log)
	appCtx.SetDB(db)
	appCtx.SetRedisClient(rc)

	// 资源回收
	defer func(appCtx *appcontext.Context) {
		if appCtx.GetDB() != nil {
			sqlDB, err := appCtx.GetDB().DB()
			if err != nil {
				panic(err)
			}

			if err := sqlDB.Close(); err != nil {
				panic(err)
			}
		}

		if appCtx.GetRedisClient() != nil {
			if err := appCtx.GetRedisClient().Close(); err != nil {
				panic(err)
			}
		}
	}(appCtx)

	// 创建 http 引擎和注册路由
	router := gin.Default()
	routes.Register(router, appCtx)

	// 调用应用钩子
	if err := httpserver.BeforeRun(router, appCtx); err != nil {
		panic(err)
	}

	// 启动 http 服务
	if err := router.Run(fmt.Sprintf("%s:%d", appCtx.GetConfig().AppConf.Host, appCtx.GetConfig().AppConf.Port)); err != nil {
		panic(err)
	}
}
