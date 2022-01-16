package script

import (
	"github.com/spf13/cobra"
	"go-scaffold/internal/app/global"
	"go.uber.org/zap"
)

// s0000000000 示例脚本
// 脚本用途：TODO
type s0000000000 struct {
	logger *zap.Logger
}

// NewS0000000000 构造函数
func NewS0000000000() *s0000000000 {
	return &s0000000000{
		logger: global.Logger(),
	}
}

// Run 脚本执行方法
func (s s0000000000) Run(cmd *cobra.Command, args []string) {
	s.logger.Sugar().Infof("%s 脚本调用成功", cmd.Use)
}
