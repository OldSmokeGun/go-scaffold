package scheduler

import (
	"github.com/robfig/cron/v3"

	"go-scaffold/internal/app/facade/cron/job"
	"go-scaffold/internal/config"
)

// Scheduler job scheduler
type Scheduler struct {
	appConf    config.App
	exampleJob *job.ExampleJob
}

// New build job scheduler
func New(
	appConf config.App,
	exampleJob *job.ExampleJob,
) *Scheduler {
	return &Scheduler{
		appConf:    appConf,
		exampleJob: exampleJob,
	}
}

// Register registers job
func (s *Scheduler) Register(server *cron.Cron) error {
	// TODO register cron job
	if _, err := server.AddFunc("*/5 * * * * *", func() {}); err != nil { // 每 5 秒钟运行一次
		return err
	}
	if _, err := server.AddJob("@daily", s.exampleJob); err != nil { // 每天 00:00 运行一次
		return err
	}
	if _, err := server.AddJob("@every 1h30m10s", s.exampleJob); err != nil { // 每 1 小时 30 分 10 秒运行一次
		return err
	}

	return nil
}
