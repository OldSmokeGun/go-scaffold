package command

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"go-scaffold/internal/config"
	"go-scaffold/pkg/ioutils"
	"go-scaffold/pkg/log"
	iklog "go-scaffold/pkg/log/kratos"
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

func (c *baseCmd) initConfig(cmd *cobra.Command, enableRemote bool) {
	configPath := cmd.Flag(flagConfig.name).Value.String()
	if configPath == "" {
		panic("config file must be specified")
	}

	cm := new(config.Config)
	configResources := []kconfig.Source{file.NewSource(configPath)}

	if enableRemote { // enable apollo
		apolloConfigEndpoint := cmd.Flag(flagApolloConfigEndpoint.name).Value.String()
		apolloConfigAppID := cmd.Flag(flagApolloConfigAppID.name).Value.String()
		apolloConfigCluster := cmd.Flag(flagApolloConfigCluster.name).Value.String()
		apolloConfigNamespace := cmd.Flag(flagApolloConfigNamespace.name).Value.String()
		apolloConfigSecret := cmd.Flag(flagApolloConfigSecret.name).Value.String()

		configResources = append(configResources, apollo.NewSource(
			apollo.WithEndpoint(apolloConfigEndpoint),
			apollo.WithAppID(apolloConfigAppID),
			apollo.WithCluster(apolloConfigCluster),
			apollo.WithNamespace(apolloConfigNamespace),
			apollo.WithSecret(apolloConfigSecret),
		))
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
	endpoint := traceConfig.Endpoint

	tracing, err := trace.New(
		endpoint,
		trace.WithServiceName(appName.String()),
		trace.WithEnv(appEnv.String()),
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
