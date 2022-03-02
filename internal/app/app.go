package app

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go-scaffold/internal/app/component"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/cron"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/pkg/migratorx"
	"go-scaffold/internal/app/transport"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	component.ProviderSet,
	cron.ProviderSet,
	transport.ProviderSet,
)

type App struct {
	logger    *log.Helper
	config    *config.Config
	db        *gorm.DB
	cron      *cron.Cron
	transport *transport.Transport
}

func New(
	logger log.Logger,
	config *config.Config,
	db *gorm.DB,
	cron *cron.Cron,
	transport *transport.Transport,
) *App {
	return &App{
		logger:    log.NewHelper(logger),
		config:    config,
		db:        db,
		cron:      cron,
		transport: transport,
	}
}

// Start 启动应用
func (a *App) Start() (err error) {
	// 数据迁移
	if a.db != nil {
		if err = migratorx.New(a.db).Run(model.MigrationTasks()); err != nil {
			return
		}

		a.logger.Info("database migration completed")
	}

	// 启动 cron 服务
	if err = a.cron.Start(); err != nil {
		return
	}

	// 启动 transport 服务
	if err = a.transport.Start(); err != nil {
		return
	}

	return nil
}

// Stop 停止应用
func (a *App) Stop(ctx context.Context) (err error) {
	// 关闭 cron 服务
	if err = a.cron.Stop(ctx); err != nil {
		return
	}

	// 关闭 transport 服务
	if err = a.transport.Stop(); err != nil {
		return
	}

	return nil
}
