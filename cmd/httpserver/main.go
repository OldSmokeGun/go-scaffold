package main

import (
	"context"
	"errors"
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
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// defaultConfigPath 默认配置文件路径
var defaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "config/httpserver.yaml")

var (
	ErrConfIncorrectValue = errors.New("in the configuration file, the value of key %s is configured incorrectly")
)

func main() {
	var (
		configPath string
		appConf    = &appconfig.Config{}
		logRotate  *rotatelogs.RotateLogs
		l          *logrus.Logger
		db         *gorm.DB
		rdb        *redis.Client
		err        error
	)

	pflag.StringVarP(&configPath, "config", "c", defaultConfigPath, "配置文件路径")
	pflag.Parse()

	// 加载配置
	if err = configurator.LoadConfig(configPath, appConf); err != nil {
		panic(err)
	}

	// 检查环境是否设置正确
	if appConf.AppConf.Env.String() != appconfig.Local.String() &&
		appConf.AppConf.Env.String() != appconfig.Test.String() &&
		appConf.AppConf.Env.String() != appconfig.Production.String() {
		panic(fmt.Sprintf(ErrConfIncorrectValue.Error(), "Env"))
	}

	// 日志轮转
	logRotate, err = rotatelogs.New(
		appConf.LogConf.Path,
		rotatelogs.WithClock(rotatelogs.Local),
	)
	defer func() {
		if err := logRotate.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}

	// 日志初始化
	appConf.LoggerConf.Output = logRotate // 设置日志的输出
	if appConf.LoggerConf != nil {
		l, err = logger.Setup(*appConf.LoggerConf)
		if err != nil {
			panic(err)
		}
	}

	// orm 初始化
	appConf.DatabaseConf.Output = logRotate // 设置日志的输出
	if appConf.DatabaseConf != nil {
		db, err = orm.Setup(*appConf.DatabaseConf)
		if err != nil {
			panic(err)
		}
	}

	// redis 初始化
	if appConf.RedisConf != nil {
		rdb, err = redisclient.Setup(*appConf.RedisConf)
		if err != nil {
			panic(err)
		}
	}

	// 创建上下文依赖
	appCtx := appcontext.New()
	appCtx.SetConfig(appConf)
	appCtx.SetLogger(l)
	appCtx.SetDB(db)
	appCtx.SetRedisClient(rdb)

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

	// 监听信号
	signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	// 初始化 router
	gin.DefaultWriter = io.MultiWriter(logRotate, os.Stdout)
	router := gin.Default()
	routes.Register(router, appCtx)

	// 调用应用钩子
	if err := httpserver.BeforeRun(router, appCtx); err != nil {
		panic(err)
	}

	// 启动 http 服务
	addr := fmt.Sprintf("%s:%d", appCtx.GetConfig().AppConf.Host, appCtx.GetConfig().AppConf.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	go func() {
		log.Printf("Listening and serving HTTP on %s", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// 等待信号
	<-signalCtx.Done()

	signalStop() // 取消信号的监听

	log.Println("The server is shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("Server forced to shutdown:", err)
	}

	log.Println("Server has been shutdown")
}
