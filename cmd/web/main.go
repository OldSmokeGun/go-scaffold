package main

import (
	"context"
	"fmt"
	"gin-scaffold/global"
	"gin-scaffold/internal/web"
	"gin-scaffold/internal/web/config"
	"gin-scaffold/pkg/configure"
	"gin-scaffold/pkg/helper"
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
var defaultConfigPath = filepath.Join(helper.RootPath(), "config/web.yaml")

func main() {
	var (
		configPath string
		conf       = &config.Config{}
		logRotate  *rotatelogs.RotateLogs
		lgr        *logrus.Logger
		db         *gorm.DB
		rdb        *redis.Client
		err        error
	)

	pflag.StringVarP(&configPath, "config", "c", defaultConfigPath, "配置文件路径")
	pflag.Parse()

	// 加载配置
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(helper.RootPath(), configPath)
	}
	configure.MustLoad(configPath, conf)

	// 检查环境是否设置正确
	switch conf.Env {
	case config.Local:
	case config.Test:
	case config.Prod:
	default:
		panic("unknown Env value: " + conf.Env)
	}

	// 日志切割
	if !filepath.IsAbs(conf.Log.Path) {
		conf.Log.Path = filepath.Join(helper.RootPath(), conf.Log.Path)
	}

	logRotate, err = rotatelogs.New(
		conf.Log.Path,
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
	if conf.Logger != nil {
		// 设置日志的路径
		if conf.Logger.Path != "" {
			if !filepath.IsAbs(conf.Logger.Path) {
				conf.Logger.Path = filepath.Join(helper.RootPath(), conf.Logger.Path)
			}
		}

		// 设置日志的输出
		conf.Logger.Output = logRotate
		lgr = logger.MustSetup(conf.Logger)
	}

	// orm 初始化
	if conf.DB != nil {
		// 设置日志的输出
		conf.DB.Output = logRotate
		db = orm.MustSetup(conf.DB)
	}

	// redis 初始化
	if conf.Redis != nil {
		rdb = redisclient.MustSetup(conf.Redis)
	}

	// 创建上下文依赖
	global.SetLogRotate(logRotate)
	global.SetConfig(conf)
	global.SetLogger(lgr)
	global.SetDB(db)
	global.SetRedisClient(rdb)

	// 资源回收
	defer func() {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				panic(err)
			}

			if err := sqlDB.Close(); err != nil {
				panic(err)
			}
		}

		if rdb != nil {
			if err := rdb.Close(); err != nil {
				panic(err)
			}
		}
	}()

	// 监听信号
	signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	// 启动 http 服务
	addr := fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: web.MustSetup(),
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
