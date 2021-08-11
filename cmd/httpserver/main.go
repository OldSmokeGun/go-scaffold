package main

import (
	"context"
	"fmt"
	"gin-scaffold/global"
	"gin-scaffold/internal/httpserver"
	"gin-scaffold/internal/httpserver/appconfig"
	"gin-scaffold/internal/httpserver/appcontext"
	"gin-scaffold/internal/httpserver/router"
	"gin-scaffold/pkg/configurator"
	"gin-scaffold/pkg/logger"
	"gin-scaffold/pkg/orm"
	"gin-scaffold/pkg/redisclient"
	"github.com/go-redis/redis/v8"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// defaultConfigPath 默认配置文件路径
var defaultConfigPath = filepath.Join(filepath.Dir(filepath.Dir(global.GetBinPath())), "config/httpserver.yaml")

var (
	ValidEnvValues         = [3]appconfig.AppEnv{appconfig.Local, appconfig.Test, appconfig.Production}
	ErrConfValueEnvInvalid = fmt.Errorf("in the configuration file,the value of key Env is invalid,allowed value is one of %s,%s,%s", appconfig.Local, appconfig.Test, appconfig.Production)
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
	configurator.MustLoadConfig(configPath, appConf)
	// 检查环境是否设置正确
	var exist bool
	for _, value := range ValidEnvValues {
		if appConf.AppConf.Env == value {
			exist = true
		}
	}
	if !exist {
		panic(ErrConfValueEnvInvalid)
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
	if appConf.LoggerConf != nil {
		// 设置日志的输出
		appConf.LoggerConf.Output = logRotate
		l = logger.MustSetup(*appConf.LoggerConf)
	}

	// orm 初始化
	if appConf.DatabaseConf != nil {
		// 设置日志的输出
		appConf.DatabaseConf.Output = logRotate
		db = orm.MustSetup(*appConf.DatabaseConf)
	}

	// redis 初始化
	if appConf.RedisConf != nil {
		rdb = redisclient.MustSetup(*appConf.RedisConf)
	}

	// 创建上下文依赖
	appCtx := appcontext.New()
	appCtx.SetLogRotate(logRotate)
	appCtx.SetConfig(appConf)
	appCtx.SetLogger(l)
	appCtx.SetDB(db)
	appCtx.SetRedisClient(rdb)

	// 资源回收
	defer func(appCtx *appcontext.Context) {
		if appCtx.DB() != nil {
			sqlDB, err := appCtx.DB().DB()
			if err != nil {
				panic(err)
			}

			if err := sqlDB.Close(); err != nil {
				panic(err)
			}
		}

		if appCtx.RedisClient() != nil {
			if err := appCtx.RedisClient().Close(); err != nil {
				panic(err)
			}
		}
	}(appCtx)

	// 监听信号
	signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	// 初始化 router
	r := router.Setup(appCtx)
	appCtx.SetRouter(r)

	// 调用应用钩子
	if err := httpserver.BeforeRun(r, appCtx); err != nil {
		panic(err)
	}

	// 启动 http 服务
	addr := fmt.Sprintf("%s:%d", appCtx.Config().AppConf.Host, appCtx.Config().AppConf.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
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
