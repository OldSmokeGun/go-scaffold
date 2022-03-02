package script

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
)

// S0000000000 示例脚本
// 脚本说明：TODO
type S0000000000 struct {
	logger *log.Helper
}

// NewS0000000000 构造函数
func NewS0000000000(logger log.Logger) *S0000000000 {
	return &S0000000000{
		logger: log.NewHelper(logger),
	}
}

// Run 脚本执行方法
func (s S0000000000) Run(cmd *cobra.Command, args []string) {
	s.logger.Infof("%s 脚本调用成功", cmd.Use)
}
