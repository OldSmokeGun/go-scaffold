package main

import (
	"context"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-scaffold/internal/app"
	appconfig "go-scaffold/internal/app/config"
	"go-scaffold/internal/app/global"
	"go-scaffold/pkg/config/local"
	"go-scaffold/pkg/config/remote"
	"go-scaffold/pkg/helper"
	"go-scaffold/pkg/helper/slicex"
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
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	rootPath   = helper.RootPath()
	logPath    string
	logLevel   string
	logFormat  string
	configPath string
	remotePath string
)

func init() {
	pflag.StringVarP(&logPath, "log.path", "", "logs/%Y%m%d.log", "日志输出路径")
	pflag.StringVarP(&logLevel, "log.level", "", "info", "日志等级（debug、info、warn、error、panic、panic、fatal）")
	pflag.StringVarP(&logFormat, "log.format", "", "json", "日志格式（text、json）")
	pflag.StringVarP(&configPath, "config", "c", filepath.Join(rootPath, "etc/app/config.yaml"), "配置文件路径")
	pflag.StringVarP(&remotePath, "config.remote", "", filepath.Join(rootPath, "etc/app/remote.yaml"), "远程配置中心配置文件路径")

	cobra.OnInitialize(setup)
}

var (
	conf         = new(appconfig.Config)  // app 配置实例
	loggerOutput *rotatelogs.RotateLogs   // 日志输出对象
	zLogger      *zap.Logger              // 日志实例
	db           *gorm.DB                 // 数据库实例
	rdb          *redis.Client            // redis 客户端实例
	tp           *sdktrace.TracerProvider // otel TracerProvider
	tracer       oteltrace.Tracer         // otel Tracer
)

func main() {
	defer cleanup()

	cmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			// 监听退出信号
			signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer signalStop()

			// 调用 app 启动钩子
			go func() {
				if err := app.Start(); err != nil {
					panic(err)
				}
			}()

			// 等待退出信号
			<-signalCtx.Done()

			signalStop() // 取消信号的监听

			zLogger.Info("the app is shutting down ...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.ShutdownWaitTime)*time.Second)
			defer cancel()

			// 关闭应用
			if err := app.Stop(ctx); err != nil {
				panic(err)
			}

			zLogger.Info("the app has been stop")
		},
	}

	// 设置全局变量
	global.SetCommand(cmd)

	// 调用 app 初始化钩子
	if err := app.Setup(); err != nil {
		panic(err)
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func setup() {
	var (
		supportedEnvs = []string{"local", "test", "prod"}
		remoteInfo    = &appconfig.Remote{} // app 远程配置中心配置信息
		err           error
	)

	if logPath != "" {
		if !filepath.IsAbs(logPath) {
			logPath = filepath.Join(helper.RootPath(), logPath)
		}

		loggerOutput, err = rotatelogs.New(
			logPath,
			rotatelogs.WithClock(rotatelogs.Local),
		)
		if err != nil {
			panic(err)
		}
	}

	// 日志初始化
	zLogger = logger.MustNew(logger.Config{
		Level:  logger.Level(logLevel),
		Format: logger.Format(logFormat),
		Output: loggerOutput,
	})

	zLogger.Info("setup resource ...")
	zLogger.Sugar().Infof("log output directory: %s", filepath.Dir(logPath))

	// 加载配置
	if configPath == "" {
		panic("config path is missing")
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(helper.RootPath(), configPath)
	}

	zLogger.Sugar().Infof("load config from: %s", configPath)

	localConfig, err := local.New(configPath)
	if err != nil {
		panic(err)
	}
	localConfig.MustLoad(conf)
	localConfig.Watch(func(v *viper.Viper, event fsnotify.Event) {
		zLogger.Sugar().Infof("config changed, filepath: %s, op: %s", event.Name, event.Op.String())

		if err := v.MergeInConfig(); err != nil {
			zLogger.Error(err.Error())
			return
		}
		if err := v.Unmarshal(conf); err != nil {
			zLogger.Error(err.Error())
			return
		}
	})

	if !conf.Priority {
		zLogger.Sugar().Infof("the local config priority is %t, will load remote config", conf.Priority)
		zLogger.Sugar().Infof("load remote info from: %s", remotePath)

		if remotePath != "" {
			if !filepath.IsAbs(remotePath) {
				remotePath = filepath.Join(helper.RootPath(), remotePath)
			}

			localRemoteConfig, err := local.New(remotePath)
			if err != nil {
				panic(err)
			}
			localRemoteConfig.MustLoad(remoteInfo)

			// 远程配置加载
			remoteConfigType := strings.TrimPrefix(path.Ext(configPath), ".")
			rc, err := remote.New(
				remoteInfo.Type,
				remoteInfo.Endpoint,
				remoteInfo.Path,
				remoteInfo.SecretKeyring,
				remoteConfigType,
				remote.WithOptions(remoteInfo.Options),
			)
			if err != nil {
				panic(err)
			}

			zLogger.Sugar().Infof("remote provider type: %s", remoteInfo.Type)
			zLogger.Sugar().Infof("remote provider endpoint: %s", remoteInfo.Endpoint)
			zLogger.Sugar().Infof("remote provider path: %s", remoteInfo.Path)
			zLogger.Sugar().Infof("remote provider config type: %s", remoteConfigType)
			zLogger.Sugar().Infof("remote provider options: %v", remoteInfo.Options)

			rc.MustLoad(conf)
			if remoteInfo.Type == "apollo" {
				rc.MustWatch(nil, func(v *viper.Viper, e interface{}) {
					event := e.(*storage.FullChangeEvent)

					zLogger.Sugar().Infof("remote config changed, namespace: %s, notification_id: %d", event.Namespace, event.NotificationID)

					for k, c := range event.Changes {
						v.Set(k, c)
					}
					if err := v.Unmarshal(conf); err != nil {
						zLogger.Error(err.Error())
						return
					}
				})
			} else {
				if err := rc.Watch(); err != nil {
					panic(err)
				}
			}
		}
	}

	// 检查环境是否设置正确
	if !slicex.InStringSlice(conf.Env, supportedEnvs) {
		panic("unsupported env value: " + conf.Env)
	}

	zLogger.Sugar().Infof("cunrrent env: %s", conf.Env)

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
			Env:         conf.Env,
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
}

// cleanup 资源回收
func cleanup() {
	zLogger.Info("cleaning resource ...")

	if rdb != nil {
		if err := rdb.Close(); err != nil {
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

	if zLogger != nil {
		if err := zLogger.Sync(); err != nil {
			log.Println(err)
		}
	}

	if loggerOutput != nil {
		if err := loggerOutput.Close(); err != nil {
			log.Println(err)
		}
	}

	if tp != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Trace.ShutdownWaitTime)*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}
}
