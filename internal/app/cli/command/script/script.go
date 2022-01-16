package script

import "github.com/spf13/cobra"

// Script 脚本接口
// 所有脚本都应该实现此接口
// 脚本名称规范为：S + 10位时间戳
type Script interface {
	Run(cmd *cobra.Command, args []string)
}
