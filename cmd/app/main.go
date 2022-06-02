package main

import (
	"context"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/fsnotify/fsnotify"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"go-scaffold/internal/app/command"
	appconfig "go-scaffold/internal/app/config"
	"go-scaffold/pkg/config/local"
	"go-scaffold/pkg/config/remote"
	"go-scaffold/pkg/log"
	"go-scaffold/pkg/path"
	"go.uber.org/zap"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	rootPath           = path.RootPath()
	logPath            string
	logLevel           string
	logFormat          string
	logCallerSkip      int    // 日志 caller 跳过层数
	configPath         string // 配置文件路径
	configRemoteType   string // 远程配置中心类型
	configRemoteFormat string // 远程配置中心类型配置格式

	apolloConfigEndpoint  string // apollo 连接地址
	apolloConfigAppID     string // apollo appID
	apolloConfigCluster   string // apollo cluster
	apolloConfigNamespace string // apollo 命名空间
	apolloConfigSecret    string // apollo 密钥

	etcdConfigEndpoint      string // etcd 连接地址
	etcdConfigPath          string // etcd 配置路径
	etcdConfigSecretKeyring string // etcd 密钥

	consulConfigEndpoint      string // consul 连接地址
	consulConfigPath          string // consul 配置路径
	consulConfigSecretKeyring string // consul 密钥
)

func init() {
	pflag.StringVarP(&logPath, "log.path", "", "logs/%Y%m%d.log", "日志输出路径")
	pflag.StringVarP(&logLevel, "log.level", "", "info", "日志等级（debug、info、warn、error、panic、panic、fatal）")
	pflag.StringVarP(&logFormat, "log.format", "", "json", "日志格式（text、json）")
	pflag.IntVarP(&logCallerSkip, "log.caller-skip", "", 4, "日志 caller 跳过层数")
	pflag.StringVarP(&configPath, "config", "f", filepath.Join(rootPath, "etc/config.yaml"), "配置文件路径")
	pflag.StringVarP(&configRemoteType, "config.remote.type", "", "", "远程配置中心类型（etcd、consul、apollo）")
	pflag.StringVarP(&configRemoteFormat, "config.remote.format", "", "", "远程配置中心类型配置格式（json、yaml）")

	pflag.StringVarP(&apolloConfigEndpoint, "config.apollo.endpoint", "", "http://localhost:8080", "apollo 连接地址")
	pflag.StringVarP(&apolloConfigAppID, "config.apollo.appid", "", "", "apollo appID")
	pflag.StringVarP(&apolloConfigCluster, "config.apollo.cluster", "", "default", "apollo cluster")
	pflag.StringVarP(&apolloConfigNamespace, "config.apollo.namespace", "", "application", "apollo 命名空间")
	pflag.StringVarP(&apolloConfigSecret, "config.apollo.secret", "", "", "apollo 密钥")

	pflag.StringVarP(&etcdConfigEndpoint, "config.etcd.endpoint", "", "http://localhost:2379", "etcd 连接地址")
	pflag.StringVarP(&etcdConfigPath, "config.etcd.path", "", "", "etcd 配置路径")
	pflag.StringVarP(&etcdConfigSecretKeyring, "config.etcd.secretKeyring", "", "", "etcd 密钥")

	pflag.StringVarP(&consulConfigEndpoint, "config.consul.endpoint", "", "http://localhost:8500", "consul 连接地址")
	pflag.StringVarP(&consulConfigPath, "config.consul.path", "", "", "consul 配置路径")
	pflag.StringVarP(&consulConfigSecretKeyring, "config.consul.secretKeyring", "", "", "consul 密钥")

	cobra.OnInitialize(setup)
}

var (
	loggerWriter *rotatelogs.RotateLogs  // 日志输出对象
	logger       *zap.Logger             // 日志实例
	configModel  = new(appconfig.Config) // app 配置实例
)

func main() {
	defer cleanup()

	cmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("starting app ...")

			// 监听退出信号
			signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer signalStop()

			appServ, appCleanup, err := initApp(loggerWriter, logger, configModel)
			if err != nil {
				panic(err)
			}
			defer appCleanup()
			// 调用 app 启动钩子
			if err := appServ.Start(signalStop); err != nil {
				panic(err)
			}

			// 等待退出信号
			<-signalCtx.Done()

			signalStop() // 取消信号的监听

			logger.Info("the app is shutting down ...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configModel.App.Timeout)*time.Second)
			defer cancel()

			// 关闭应用
			if err := appServ.Stop(ctx); err != nil {
				panic(err)
			}

			logger.Info("the app has been stop")
		},
	}

	command.Setup(cmd, func() (*command.Command, func(), error) {
		return initCommand(loggerWriter, logger, configModel)
	})

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func setup() {
	var err error

	if logPath != "" {
		if !filepath.IsAbs(logPath) {
			logPath = filepath.Join(path.RootPath(), logPath)
		}

		loggerWriter, err = rotatelogs.New(
			logPath,
			rotatelogs.WithClock(rotatelogs.Local),
		)
		if err != nil {
			panic(err)
		}
	}

	// 日志初始化
	var writer io.Writer
	if loggerWriter == nil {
		writer = os.Stdout
	} else {
		writer = io.MultiWriter(os.Stdout, loggerWriter)
	}
	logger = log.New(
		log.WithLevel(log.Level(logLevel)),
		log.WithFormat(log.Format(logFormat)),
		log.WithWriter(writer),
		log.WithCallerSkip(logCallerSkip),
	)

	logger.Info("setup resource ...")
	logger.Sugar().Infof("log output directory: %s", filepath.Dir(logPath))

	// 加载配置
	if configPath == "" {
		panic("config path is missing")
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(path.RootPath(), configPath)
	}

	logger.Sugar().Infof("load config from: %s", configPath)

	localConfig, err := local.New(configPath)
	if err != nil {
		panic(err)
	}
	localConfig.MustLoad(configModel)
	localConfig.Watch(func(v *viper.Viper, event fsnotify.Event) {
		logger.Sugar().Infof("config changed, filepath: %s, op: %s", event.Name, event.Op.String())

		if err := v.MergeInConfig(); err != nil {
			logger.Error(err.Error())
			return
		}
		if err := v.Unmarshal(configModel); err != nil {
			logger.Error(err.Error())
			return
		}
	})

	if configRemoteType != "" {
		logger.Sugar().Infof("enable remote config, config will be loaded from %s", configRemoteType)

		var (
			rc            *remote.Remote
			rcEndpoint    string
			rcPath        string
			rcSecret      string
			rcOptions     = map[string]interface{}{}
			rcOptionsFunc remote.OptionFunc
		)

		switch configRemoteType {
		case "etcd":
			rcEndpoint = etcdConfigEndpoint
			rcPath = etcdConfigPath
			rcSecret = etcdConfigSecretKeyring
		case "consul":
			rcEndpoint = consulConfigEndpoint
			rcPath = consulConfigPath
			rcSecret = consulConfigSecretKeyring
		case "apollo":
			rcEndpoint = apolloConfigEndpoint
			rcPath = apolloConfigNamespace
			rcSecret = apolloConfigSecret
			rcOptions = map[string]interface{}{
				"AppID":   apolloConfigAppID,
				"Cluster": apolloConfigCluster,
			}
			rcOptionsFunc = remote.WithOptions(rcOptions)
		}
		rc, err = remote.New(configRemoteType, rcEndpoint, rcPath, rcSecret, configRemoteFormat, rcOptionsFunc)
		if err != nil {
			panic(err)
		}

		logger.Sugar().Infof("remote provider type: %s", configRemoteType)
		logger.Sugar().Infof("remote provider endpoint: %s", rcEndpoint)
		logger.Sugar().Infof("remote provider path: %s", rcPath)
		logger.Sugar().Infof("remote provider config format: %s", configRemoteFormat)
		logger.Sugar().Infof("remote provider options: %v", rcOptions)

		rc.MustLoad(configModel)
		if configRemoteType == "apollo" {
			rc.MustWatch(nil, func(v *viper.Viper, e interface{}) {
				event := e.(*storage.FullChangeEvent)

				logger.Sugar().Infof("remote config changed, namespace: %s, notification_id: %d", event.Namespace, event.NotificationID)

				for k, c := range event.Changes {
					v.Set(k, c)
				}
				if err := v.Unmarshal(configModel); err != nil {
					logger.Error(err.Error())
					return
				}
			})
		} else {
			if err := rc.Watch(); err != nil {
				panic(err)
			}
		}

	}

	if err := appconfig.Loaded(logger, configModel); err != nil {
		panic(err)
	}

	// 检查环境是否设置正确
	if !funk.ContainsString(appconfig.SupportedEnvs, configModel.App.Env.String()) {
		panic("unsupported env value: " + configModel.App.Env)
	}

	logger.Sugar().Infof("cunrrent env: %s", configModel.App.Env)
}

// cleanup 资源回收
func cleanup() {
	if logger != nil {
		logger.Info("resource cleaning ...")
	}

	if loggerWriter != nil {
		if err := loggerWriter.Close(); err != nil {
			panic(err)
		}
	}

	if logger != nil {
		if err := logger.Sync(); err != nil {
			logger.Error(err.Error())
		}
	}
}
