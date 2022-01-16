package job

import (
	"go-scaffold/internal/app/global"
	"go.uber.org/zap"
)

// exampleJob 示例任务
// 任务说明：TODO
type exampleJob struct {
	logger *zap.Logger
}

// NewExampleJob 构造函数
func NewExampleJob() *exampleJob {
	return &exampleJob{
		logger: global.Logger(),
	}
}

// Run 任务执行方法
func (s exampleJob) Run() {
	s.logger.Info("exampleJob 任务执行成功")
}
