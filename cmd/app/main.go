package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-scaffold/internal/app"
	appconfig "go-scaffold/internal/app/config"
	"go-scaffold/internal/app/global"
	"go-scaffold/pkg/config"
	"go-scaffold/pkg/helper"
	"go-scaffold/pkg/logger"
	"go-scaffold/pkg/orm"
	"go-scaffold/pkg/redisclient"
	"go-scaffold/pkg/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// defaultConfigPath 默认配置文件路径
var defaultConfigPath = filepath.Join(helper.RootPath(), "etc/app.yaml")

func main() {
	var (
		configPath   string
		conf         = &appconfig.Config{}    // app 配置实例
		loggerOutput *rotatelogs.RotateLogs   // 日志输出对象
		zLogger      *zap.Logger              // 日志实例
		db           *gorm.DB                 // 数据库实例
		rdb          *redis.Client            // redis 客户端实例
		tp           *sdktrace.TracerProvider // otel TracerProvider
		tracer       oteltrace.Tracer         // otel Tracer
		err          error
	)

	pflag.StringVarP(&configPath, "etc", "c", defaultConfigPath, "配置文件路径")

	// 加载配置
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(helper.RootPath(), configPath)
	}
	config.MustLoad(configPath, conf, func(v *viper.Viper, model interface{}, e fsnotify.Event) {
		if err := v.MergeInConfig(); err != nil {
			panic(err)
		}
		if err := v.Unmarshal(model); err != nil {
			panic(err)
		}
	})

	// 检查环境是否设置正确
	switch conf.Env {
	case appconfig.Local:
	case appconfig.Test:
	case appconfig.Prod:
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

	// tracer 初始化
	if conf.Trace != nil {
		tp = trace.MustNew(trace.Config{
			Endpoint:    conf.Trace.Endpoint,
			ServiceName: conf.Name,
			Env:         conf.Env.String(),
		})
		tracer = tp.Tracer(conf.Name)
	}

	// 设置全局变量
	global.SetLoggerOutput(loggerOutput)
	global.SetConfig(conf)
	global.SetLogger(zLogger)
	global.SetDB(db)
	global.SetRedisClient(rdb)
	global.SetTracer(tracer)

	// 资源回收
	defer func() {
		if rdb != nil {
			if err = rdb.Close(); err != nil {
				log.Println(err.Error())
			}
		}

		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Println(err.Error())
			}

			if err := sqlDB.Close(); err != nil {
				log.Println(err.Error())
			}
		}

		if err = zLogger.Sync(); err != nil {
			log.Println(err)
		}

		if err = loggerOutput.Close(); err != nil {
			log.Println(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Trace.ShutdownWaitTime)*time.Second)
		defer cancel()
		if err = tp.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	cmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			// 监听退出信号
			signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer signalStop()

			// 调用 app 启动钩子
			go func() {
				if err = app.Start(); err != nil {
					panic(err)
				}
			}()

			// 等待退出信号
			<-signalCtx.Done()

			signalStop() // 取消信号的监听

			log.Println("the app is shutting down ...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.ShutdownWaitTime)*time.Second)
			defer cancel()

			// 关闭应用
			if err = app.Stop(ctx); err != nil {
				panic(err)
			}

			log.Println("the app has been stop")
		},
	}

	// 设置全局变量
	global.SetCommand(cmd)

	// 调用 app 初始化钩子
	if err = app.Setup(); err != nil {
		panic(err)
	}

	if err = cmd.Execute(); err != nil {
		panic(err)
	}
}
