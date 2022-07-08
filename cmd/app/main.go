package main

import (
	"context"
	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
	"go-scaffold/internal/app/command"
	appconfig "go-scaffold/internal/app/config"
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
	appName     = "go-scaffold"
	hostname, _ = os.Hostname()
)

var (
	rootPath              = path.RootPath()
	logPath               string // log output path
	logLevel              string // log level
	logFormat             string // log output format
	logCallerSkip         int    // log caller skip
	configPath            string // configuration file path
	apolloConfigEnable    bool   // enable apollo
	apolloConfigEndpoint  string // apollo endpoint
	apolloConfigAppID     string // apollo appID
	apolloConfigCluster   string // apollo cluster
	apolloConfigNamespace string // apollo namespace
	apolloConfigSecret    string // apollo secret
)

func init() {
	pflag.StringVarP(&logPath, "log.path", "", "logs/%Y%m%d.log", "log output path")
	pflag.StringVarP(&logLevel, "log.level", "", "info", "log level（debug、info、warn、error、panic、fatal）")
	pflag.StringVarP(&logFormat, "log.format", "", "json", "log output format（text、json）")
	pflag.IntVarP(&logCallerSkip, "log.caller-skip", "", 4, "log caller skip")
	pflag.StringVarP(&configPath, "config", "f", filepath.Join(rootPath, "etc/config.yaml"), "configuration file path")
	pflag.BoolVarP(&apolloConfigEnable, "config.apollo.enable", "", false, "enable apollo")
	pflag.StringVarP(&apolloConfigEndpoint, "config.apollo.endpoint", "", "http://localhost:8080", "apollo endpoint")
	pflag.StringVarP(&apolloConfigAppID, "config.apollo.appid", "", "", "apollo appID")
	pflag.StringVarP(&apolloConfigCluster, "config.apollo.cluster", "", "default", "apollo cluster")
	pflag.StringVarP(&apolloConfigNamespace, "config.apollo.namespace", "", "application", "apollo namespace")
	pflag.StringVarP(&apolloConfigSecret, "config.apollo.secret", "", "", "apollo secret")

	cobra.OnInitialize(setup)
}

var (
	loggerWriter *rotatelogs.RotateLogs  // log writer
	logger       klog.Logger             // kratos log interface
	hLogger      *klog.Helper            // kratos log helper
	zLogger      *zap.Logger             // zap logger
	config       kconfig.Config          // kratos config interface
	configModel  = new(appconfig.Config) // app config model
)

func main() {
	defer cleanup()

	cmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			hLogger.Info("starting app ...")

			appServ, appCleanup, err := initApp(loggerWriter, logger, zLogger, configModel)
			if err != nil {
				panic(err)
			}
			defer appCleanup()

			// app start
			appStop, err := appServ.Start()
			if err != nil {
				panic(err)
			}

			// monitor signal
			signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer signalStop()

			// waiting exit signal ...
			select {
			case err = <-appStop:
				hLogger.Error(err)
			case <-signalCtx.Done():
			}
			signalStop()

			hLogger.Info("the app is shutting down ...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configModel.App.Timeout)*time.Second)
			defer cancel()

			// stop app
			appServ.Stop(ctx)
		},
	}

	command.Setup(cmd, func() (*command.Command, func(), error) {
		return initCommand(loggerWriter, logger, zLogger, configModel)
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

	// initialization logger
	var writer io.Writer
	if loggerWriter == nil {
		writer = os.Stdout
	} else {
		writer = io.MultiWriter(os.Stdout, loggerWriter)
	}
	zLogger = log.New(
		log.WithLevel(log.Level(logLevel)),
		log.WithFormat(log.Format(logFormat)),
		log.WithWriter(writer),
		log.WithCallerSkip(logCallerSkip),
	)
	logger = klog.With(
		kzap.NewLogger(zLogger),
		"service.id", hostname,
		"service.name", appName,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	hLogger = klog.NewHelper(logger)

	hLogger.Info("initializing resource ...")
	hLogger.Infof("the log output directory: %s", filepath.Dir(logPath))

	// load configuration
	if configPath == "" {
		panic("config path is missing")
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(path.RootPath(), configPath)
	}

	hLogger.Infof("load config from: %s", configPath)

	configResources := []kconfig.Source{file.NewSource(configPath)}
	if apolloConfigEnable { // enable apollo
		hLogger.Infof("enable remote config, config will be loaded from remote config center")

		configResources = append(configResources, apollo.NewSource(
			apollo.WithEndpoint(apolloConfigEndpoint),
			apollo.WithAppID(apolloConfigAppID),
			apollo.WithCluster(apolloConfigCluster),
			apollo.WithNamespace(apolloConfigNamespace),
			apollo.WithSecret(apolloConfigSecret),
		))
	}

	config = kconfig.New(
		kconfig.WithSource(configResources...),
		kconfig.WithLogger(logger),
	)

	if err := config.Load(); err != nil {
		panic(err)
	}

	if err := config.Scan(configModel); err != nil {
		panic(err)
	}

	if err := appconfig.Loaded(logger, config, configModel); err != nil {
		panic(err)
	}

	// check that the environment is set correctly
	if !funk.ContainsString(appconfig.SupportedEnvs, configModel.App.Env.String()) {
		panic("unsupported env value: " + configModel.App.Env)
	}

	hLogger.Infof("current env: %s", configModel.App.Env)
}

// cleanup recycle resources
func cleanup() {
	if hLogger != nil {
		hLogger.Info("resource cleaning ...")
	}

	if config != nil {
		if err := config.Close(); err != nil {
			hLogger.Error(err.Error())
		}
	}

	if loggerWriter != nil {
		if err := loggerWriter.Close(); err != nil {
			panic(err)
		}
	}

	if zLogger != nil {
		if err := zLogger.Sync(); err != nil {
			hLogger.Error(err.Error())
		}
	}
}
