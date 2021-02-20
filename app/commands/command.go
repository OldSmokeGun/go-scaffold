package commands

import (
	"github.com/spf13/cobra"
)

// Register 注册根命令的 flag 或子命令
func Register(cmd *cobra.Command)  {
	// 注册根命令的 flag
	// ...

	// 注册子命令
	cmd.AddCommand(&cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			// database migrate ...
		},
	})
}
