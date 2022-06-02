package job

import (
	"go.uber.org/zap"
)

// Example 示例任务
// 任务说明：TODO
type Example struct {
	logger *zap.Logger
}

// NewExample 构造函数
func NewExample(
	logger *zap.Logger,
) *Example {
	return &Example{
		logger: logger,
	}
}

// Run 任务执行方法
func (s Example) Run() {
	s.logger.Info("Example 任务执行成功")
}
