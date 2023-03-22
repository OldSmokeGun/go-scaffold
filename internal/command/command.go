package command

import (
	"github.com/spf13/cobra"
)

type cmder interface {
	getCmd() *cobra.Command
	addCommands(commands ...cmder)
}

// NewCommand return command
func NewCommand() *cobra.Command {
	return newRootCmd().getCmd()
}
