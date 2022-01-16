package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
)

// cronServer cron 服务
var cronServer *cron.Cron

// Start cron 服务启动
func Start() (err error) {
	cronServer = cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.PrintfLogger(log.Default())),
			cron.DelayIfStillRunning(cron.PrintfLogger(log.Default())),
		),
	)

	// TODO 编写 cron 任务
	// if _, err = cronServer.AddFunc("*/5 * * * * *", func() {}); err != nil { // 每 5 秒钟运行一次
	// 	return err
	// }
	// if _, err = cronServer.AddJob("@daily", job.NewExampleJob()); err != nil { // 每天 00:00 运行一次
	// 	return err
	// }
	// if _, err = cronServer.AddJob("@every 1h30m10s", job.NewExampleJob()); err != nil { // 每 1 小时 30 分 10 秒运行一次
	// 	return err
	// }

	// 启动 cron 服务
	cronServer.Start()

	return nil
}

// Stop cron 服务关闭
func Stop(ctx context.Context) (err error) {
	cronServer.Stop()

	return
}
