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
	"github.com/spf13/pflag"
	"go.uber.org/zap"
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
		configPath   string
		conf         = &config.Config{}
		loggerOutput *rotatelogs.RotateLogs
		zLogger      *zap.Logger
		db           *gorm.DB
		rdb          *redis.Client
		err          error
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

	if conf.Log.Path != "" {
		logPath := conf.Log.Path
		if !filepath.IsAbs(conf.Log.Path) {
			logPath = filepath.Join(helper.RootPath(), conf.Log.Path)
		}

		// 日志切割
		loggerOutput, err = rotatelogs.New(
			logPath,
			rotatelogs.WithClock(rotatelogs.Local),
		)
		if err != nil {
			panic(err)
		}
		defer loggerOutput.Close()
		// 日志初始化
		zLogger = logger.MustNew(logger.Config{
			Path:   conf.Log.Path,
			Level:  conf.Log.Level,
			Format: conf.Log.Format,
			Caller: conf.Log.Caller,
			Mode:   conf.Log.Mode,
			Output: loggerOutput,
		})
	}

	// orm 初始化
	if conf.DB != nil {
		db = orm.MustNew(*conf.DB)
	}

	// redis 初始化
	if conf.Redis != nil {
		rdb = redisclient.MustNew(*conf.Redis)
	}

	// 创建上下文依赖
	global.SetLoggerOutput(loggerOutput)
	global.SetConfig(conf)
	global.SetLogger(zLogger)
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
		log.Printf("Http server started on %s", addr)

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
