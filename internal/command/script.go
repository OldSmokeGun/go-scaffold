package command

import "github.com/spf13/cobra"

type scriptCmd struct {
	*baseCmd
}

func newScriptCmd() *scriptCmd {
	c := &scriptCmd{new(baseCmd)}

	c.cmd = &cobra.Command{
		Use:   "script",
		Short: "some business script",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				panic(err)
			}
		},
	}

	addRemoteConfigFlag(c.cmd, false)
	addLoggerFlag(c.cmd, true)

	c.addCommands(
		newExampleCmd(),
	)

	return c
}

type exampleCmd struct {
	*baseCmd
}

func newExampleCmd() *exampleCmd {
	c := &exampleCmd{new(baseCmd)}

	c.cmd = &cobra.Command{
		Use:   "example",
		Short: "example script",
		Run: func(cmd *cobra.Command, args []string) {
			c.initRuntime(cmd)
			c.initLogger(cmd)
			defer c.closeLogger()

			c.initConfig(cmd)
			defer c.closeConfig()

			c.run(cmd)
		},
	}

	return c
}

func (c *exampleCmd) run(cmd *cobra.Command) {
	script, cleanup, err := newExampleScript(cmd.Context(), c.appName, c.appEnv, c.logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := script.Run(cmd); err != nil {
		panic(err)
	}
}
