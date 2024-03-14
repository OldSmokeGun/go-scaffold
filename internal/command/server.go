package command

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type serverCmd struct {
	*baseCmd
}

func newServerCmd() *serverCmd {
	c := &serverCmd{new(baseCmd)}

	c.cmd = &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve"},
		Short:   "http and grpc server",
		Run: func(cmd *cobra.Command, args []string) {
			c.initRuntime(cmd)
			c.initLogger(cmd)
			defer c.closeLogger()

			c.initConfig(cmd)
			defer c.closeConfig()

			c.initTrace(cmd)
			defer c.closeTrace(cmd.Context())

			c.watchConfig()

			c.run(cmd.Context())
		},
	}

	addRemoteConfigFlag(c.cmd, false)
	addLoggerFlag(c.cmd, false)

	return c
}

func (c *serverCmd) run(ctx context.Context) {
	stop := make(chan error, 1)

	server, cleanup, err := initServer(ctx, c.appName, c.appEnv, c.logger, c.trace)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	go func() {
		if err := server.Start(); err != nil {
			stop <- err
		}
	}()

	signalCtx, signalStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	select {
	case err := <-stop:
		c.logger.Error("start server error", slog.Any("error", err))
		return
	case <-signalCtx.Done():
	}

	if err := server.Stop(); err != nil {
		panic(err)
	}
}
