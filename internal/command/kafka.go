package command

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type kafkaCmd struct {
	*baseCmd
}

func newKafkaCmd() *kafkaCmd {
	c := &kafkaCmd{new(baseCmd)}

	c.cmd = &cobra.Command{
		Use:   "kafka",
		Short: "kafka consumer",
		Run: func(cmd *cobra.Command, args []string) {
			c.initRuntime(cmd)
			c.initLogger(cmd)
			defer c.closeLogger()

			apolloConfigEnable, err := cmd.Flags().GetBool(flagApolloConfigEnable.name)
			if err != nil {
				panic(err)
			}

			c.initConfig(cmd, apolloConfigEnable)
			defer c.closeConfig()

			c.initTrace(cmd)
			defer c.closeTrace(cmd.Context())

			c.watchConfig()

			c.run(cmd.Context())
		},
	}

	addApolloConfigFlag(c.cmd, false)
	addLoggerFlag(c.cmd, false)

	return c
}

func (c *kafkaCmd) run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	kafka, cleanup, err := initKafka(ctx, c.appName, c.appEnv, c.logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	kafka.Start(ctx)

	signalCtx, signalStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	select {
	case <-signalCtx.Done():
		cancel()
	}
}
