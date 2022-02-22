package app

import (
	"context"
	"go-scaffold/internal/app/cli"
	"go-scaffold/internal/app/cron"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/pkg/migratorx"
	"go-scaffold/internal/app/rest"
)

// Setup 应用初始化钩子
// 在这里可以执行一些初始化操作，例如：命令行 flag 的注册
func Setup() (err error) {
	// 初始化命令行
	if err = cli.Setup(); err != nil {
		return
	}

	return nil
}

// Start 应用启动钩子
func Start() (err error) {
	// 数据迁移
	if global.DB() != nil {
		if err = migratorx.New(global.DB()).Run(model.MigrationTasks()); err != nil {
			return
		}

		global.Logger().Info("database migration completed")
	}

	// 启动 cron 服务
	if err = cron.Start(); err != nil {
		return
	}

	// 启动 HTTP 服务
	if err = rest.Start(); err != nil {
		return
	}

	return nil
}

// Stop 应用关闭钩子
func Stop(ctx context.Context) (err error) {
	// 关闭 cron 服务
	if err = cron.Stop(ctx); err != nil {
		return
	}

	// 关闭 HTTP 服务
	if err = rest.Stop(ctx); err != nil {
		return
	}

	return nil
}
