package greet

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (h handler) Default(cmd *cobra.Command, args []string) {
	h.logger.Infof("%s 命令调用成功", cmd.Use)
	fmt.Println("Hello")
}
