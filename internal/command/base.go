package command

import (
	"context"
	"io"
	"log/slog"
	"os"

	remote "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cobra"
	etcdctl "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"

	"go-scaffold/internal/config"
	"go-scaffold/pkg/ioutils"
	"go-scaffold/pkg/log"
	iklog "go-scaffold/pkg/log/kratos"
	otlog "go-scaffold/pkg/log/otel"
	"go-scaffold/pkg/trace"
)

type baseCmd struct {
	cmd *cobra.Command

	appName config.AppName
	appEnv  config.Env

	config       *config.Config
	configSource kconfig.Config

	logger       *slog.Logger
	loggerWriter io.WriteCloser

	trace *trace.Trace
}

func (c *baseCmd) getCmd() *cobra.Command {
	return c.cmd
}

func (c *baseCmd) addCommands(commands ...cmder) {
	for _, command := range commands {
		c.cmd.AddCommand(command.getCmd())
	}
}

func (c *baseCmd) initRuntime(cmd *cobra.Command) {
	c.appName = getAppName(cmd)
	appEnv := getAppEnvironment(cmd)
	if !appEnv.Check() {
		panic("unsupported environment: " + appEnv)
	}
	c.appEnv = appEnv
}

func (c *baseCmd) initLogger(cmd *cobra.Command) {
	var (
		logger       *slog.Logger
		loggerWriter io.WriteCloser
	)

	appName := getAppName(cmd)
	appEnv := getAppEnvironment(cmd)
	logPath := cmd.Flag(flagLoggerPath.name).Value.String()
	logLevel := cmd.Flag(flagLoggerLevel.name).Value.String()
	logFormat := cmd.Flag(flagLoggerFormat.name).Value.String()

	loggerWriter = os.Stdout
	if logPath != "" {
		rotate, err := rotatelogs.New(logPath, rotatelogs.WithClock(rotatelogs.Local))
		if err != nil {
			panic(err)
		}

		loggerWriter = ioutils.MultiWriteCloser(os.Stdout, rotate)
	}

	hostname, _ := os.Hostname()
	logger = log.New(
		log.WithLevel(log.Level(logLevel)),
		log.WithFormat(log.Format(logFormat)),
		log.WithWriter(loggerWriter),
		log.WithAttrs([]slog.Attr{
			slog.String("host", hostname),
			slog.String("app", appName.String()),
			slog.String("env", appEnv.String()),
		}),
	)

	klog.SetLogger(iklog.NewLogger(logger))

	c.logger = logger
	c.loggerWriter = loggerWriter
}

func (c *baseCmd) initConfig(cmd *cobra.Command) {
	configPath := cmd.Flag(flagConfig.name).Value.String()
	if configPath == "" {
		panic("config file must be specified")
	}

	cm := new(config.Config)
	configResources := []kconfig.Source{file.NewSource(configPath)}

	if cmd.Flags().Lookup(flagRemoteConfigEnable.name) != nil { // enable remote config
		enableRemote, err := cmd.Flags().GetBool(flagRemoteConfigEnable.name)
		if err != nil {
			panic(err)
		}

		if enableRemote {
			remoteConfigEndpoints, err := cmd.Flags().GetStringSlice(flagRemoteConfigEndpoints.name)
			if err != nil {
				panic(err)
			}
			remoteConfigTimeout, err := cmd.Flags().GetDuration(flagRemoteConfigTimeout.name)
			if err != nil {
				panic(err)
			}
			remoteConfigPathPrefix, err := cmd.Flags().GetString(flagRemoteConfigPathPrefix.name)
			if err != nil {
				panic(err)
			}

			client, err := etcdctl.New(etcdctl.Config{
				Endpoints:   remoteConfigEndpoints,
				DialTimeout: remoteConfigTimeout,
				DialOptions: []grpc.DialOption{grpc.WithBlock()},
			})
			if err != nil {
				panic(err)
			}
			remoteSource, err := remote.New(
				client,
				remote.WithContext(cmd.Context()),
				remote.WithPath(remoteConfigPathPrefix),
				remote.WithPrefix(true),
			)
			if err != nil {
				panic(err)
			}
			configResources = append(configResources, remoteSource)
		}
	}

	src := kconfig.New(kconfig.WithSource(configResources...))
	if err := src.Load(); err != nil {
		panic(err)
	}
	if err := src.Scan(cm); err != nil {
		panic(err)
	}

	config.SetConfig(cm)

	c.config = cm
	c.configSource = src
}

func (c *baseCmd) initTrace(cmd *cobra.Command) {
	c.mustConfig()

	traceConfig, err := config.GetTrace()
	if config.IsNotConfigured(err) {
		return
	} else if err != nil {
		panic(err)
	}

	appName := getAppName(cmd)
	appEnv := getAppEnvironment(cmd)

	tracing, err := trace.New(
		cmd.Context(),
		trace.OTPLProtocol(traceConfig.Protocol),
		traceConfig.Endpoint,
		trace.WithServiceName(appName.String()),
		trace.WithEnv(appEnv.String()),
		trace.WithErrorLogger(otlog.NewLogger(c.logger)),
	)
	if err != nil {
		panic(err)
	}

	otel.SetTracerProvider(tracing.TracerProvider())
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	c.trace = tracing
	c.logger = log.NewWithHandler(otlog.NewHandler(c.logger.Handler()))
}

func (c *baseCmd) closeConfig() {
	if c.configSource != nil {
		if err := c.configSource.Close(); err != nil {
			panic(err)
		}
	}
}

func (c *baseCmd) closeLogger() {
	if c.loggerWriter != nil {
		if err := c.loggerWriter.Close(); err != nil {
			panic(err)
		}
	}
}

func (c *baseCmd) closeTrace(ctx context.Context) {
	if c.trace != nil {
		if err := c.trace.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

func (c *baseCmd) mustLogger() {
	if c.logger == nil || c.loggerWriter == nil {
		panic("please initialize logger first")
	}
}

func (c *baseCmd) mustConfig() {
	if c.configSource == nil || !config.HasConfigured() {
		panic("please load the configuration first")
	}
}

// watchConfig listen for configuration key changes
func (c *baseCmd) watchConfig() {
	c.mustLogger()
	c.mustConfig()

	if err := config.Watch(c.logger, c.configSource, c.config); err != nil {
		panic(err)
	}
}
