package greet

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (h handler) To(cmd *cobra.Command, args []string) {
	h.logger.Infof("%s 命令调用成功", cmd.Use)
	fmt.Printf("Hello, %s\n", args[0])
}
