package cron

import (
	"context"
	"log/slog"

	"github.com/google/wire"
	"github.com/robfig/cron/v3"

	"go-scaffold/internal/app/facade/cron/job"
	"go-scaffold/internal/app/facade/cron/scheduler"
	clog "go-scaffold/pkg/log/cron"
)

var ProviderSet = wire.NewSet(
	// cron job
	job.NewExampleJob,
	// scheduler
	scheduler.New,
	// cron server
	New,
)

// Cron crontab server
type Cron struct {
	server    *cron.Cron
	logger    *slog.Logger
	scheduler *scheduler.Scheduler
}

// New build crontab server
func New(
	logger *slog.Logger,
	scheduler *scheduler.Scheduler,
) (*Cron, error) {
	cLogger := clog.NewLogger(logger, false)
	server := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(cLogger),
		cron.WithChain(
			cron.Recover(cLogger),
			cron.DelayIfStillRunning(cLogger),
		),
	)

	return &Cron{
		server:    server,
		logger:    logger,
		scheduler: scheduler,
	}, nil
}

// Start cron server
func (c *Cron) Start() error {
	if err := c.scheduler.Register(c.server); err != nil {
		return err
	}

	// start cron server
	c.server.Start()

	return nil
}

// Stop cron server
func (c *Cron) Stop() context.Context {
	return c.server.Stop()
}
