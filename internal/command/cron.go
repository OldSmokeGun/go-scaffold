package command

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type cronCmd struct {
	*baseCmd
}

func newCronCmd() *cronCmd {
	c := &cronCmd{new(baseCmd)}

	c.cmd = &cobra.Command{
		Use:   "cron",
		Short: "cron server",
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

func (c *cronCmd) run(ctx context.Context) {
	c.logger.Info("cron server starting...")

	cron, cleanup, err := initCron(ctx, c.appName, c.appEnv, c.logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := cron.Start(); err != nil {
		panic(err)
	}

	signalCtx, signalStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	<-signalCtx.Done()

	cronCtx := cron.Stop()
	if err := cronCtx.Err(); err != nil {
		panic(err)
	}
	<-cronCtx.Done()

	c.logger.Info("the cron server stopped")
}
