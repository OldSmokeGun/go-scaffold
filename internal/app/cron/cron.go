package cron

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"go-scaffold/internal/app/cron/job"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

var ProviderSet = wire.NewSet(job.NewExample, New)

type Cron struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.Client
	server *cron.Cron

	exampleJob *job.Example
}

func New(
	logger *zap.Logger,
	db *gorm.DB,
	rdb *redis.Client,
	exampleJob *job.Example,
) (*Cron, error) {
	server := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.PrintfLogger(log.Default())),
			cron.DelayIfStillRunning(cron.PrintfLogger(log.Default())),
		),
	)

	return &Cron{
		logger:     logger,
		db:         db,
		rdb:        rdb,
		server:     server,
		exampleJob: exampleJob,
	}, nil
}

// Start cron 服务启动
func (c *Cron) Start() (err error) {
	// TODO 编写 cron 任务
	// if _, err = c.server.AddFunc("*/5 * * * * *", func() {}); err != nil { // 每 5 秒钟运行一次
	// 	return err
	// }
	// if _, err = c.server.AddJob("@daily", c.exampleJob); err != nil { // 每天 00:00 运行一次
	// 	return err
	// }
	// if _, err = c.server.AddJob("@every 1h30m10s", c.exampleJob); err != nil { // 每 1 小时 30 分 10 秒运行一次
	// 	return err
	// }

	// 启动 cron 服务
	c.server.Start()

	c.logger.Info("cron server started")
	return nil
}

// Stop cron 服务关闭
func (c *Cron) Stop(ctx context.Context) (err error) {
	c.server.Stop()

	c.logger.Info("cron server has been stop")
	return nil
}
