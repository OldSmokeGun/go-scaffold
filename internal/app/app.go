package app

import (
	"context"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/pkg/migratorx"
	"go-scaffold/internal/app/rest"
)

// Start 应用启动入口
func Start() (err error) {
	// 数据迁移
	if global.DB() != nil {
		if err = migratorx.New(global.DB()).Exec(model.MigrationTasks()); err != nil {
			return
		}
	}

	// 启动 HTTP 接口服务
	if err = rest.Start(); err != nil {
		return
	}

	return nil
}

// Stop 应用关闭入口
func Stop(ctx context.Context) (err error) {
	if err = rest.Stop(ctx); err != nil {
		return
	}

	return nil
}
