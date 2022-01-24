package model

import (
	"go-scaffold/internal/app/pkg/migratorx"
	"gorm.io/plugin/soft_delete"
)

// BaseModel 基础模型
// 自动更新时间戳、软删除
type BaseModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt int64
	UpdatedAt int64
	DeletedAt soft_delete.DeletedAt `gorm:"index"`
}

// MigrationTasks 注册数据迁移任务
func MigrationTasks() []*migratorx.Task {
	return []*migratorx.Task{
		{Comment: "用户表", Model: &User{}},
	}
}
