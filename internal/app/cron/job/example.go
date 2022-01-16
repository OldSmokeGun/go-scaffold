package job

import (
	"go-scaffold/internal/app/global"
	"go.uber.org/zap"
)

// example 示例任务
// 任务说明：TODO
type example struct {
	logger *zap.Logger
}

// NewExample 构造函数
func NewExample() *example {
	return &example{
		logger: global.Logger(),
	}
}

// Run 任务执行方法
func (s example) Run() {
	s.logger.Info("example 任务执行成功")
}
