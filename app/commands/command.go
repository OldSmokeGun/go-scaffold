package commands

import (
	"github.com/spf13/cobra"
)

func Register(cmd *cobra.Command)  {
	cmd.AddCommand(&cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			// database migrate ...
		},
	})
}
