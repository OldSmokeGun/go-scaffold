package script

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

var ProviderSet = wire.NewSet(NewS0000000000)

// Script 脚本接口
// 所有脚本都应该实现此接口
// 脚本名称规范为：S + 10位时间戳
type Script interface {
	Run(cmd *cobra.Command, args []string)
}
