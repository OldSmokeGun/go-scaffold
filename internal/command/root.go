package command

import (
	"github.com/spf13/cobra"
)

type rootCmd struct {
	*baseCmd
}

func newRootCmd() *rootCmd {
	c := &rootCmd{new(baseCmd)}

	c.cmd = &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				panic(err)
			}
		},
	}

	addConfigFlag(c.cmd, true)
	addAppRuntimeFlag(c.cmd, true)

	c.addCommands(
		newServerCmd(),
		newCronCmd(),
		newMigrateCmd(),
		newKafkaCmd(),
		newScriptCmd(),
	)

	return c
}
