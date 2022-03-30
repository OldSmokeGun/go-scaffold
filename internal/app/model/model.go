package model

import (
	"go-scaffold/internal/app/pkg/migratorx"
	"gorm.io/plugin/soft_delete"
)

// BaseModel 基础模型
// 自动更新时间戳、软删除
type BaseModel struct {
	Id        uint64                `gorm:"primaryKey"`
	CreatedAt int64                 `gorm:"NOT NULL;DEFAULT:0"`
	UpdatedAt int64                 `gorm:"NOT NULL;DEFAULT:0"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;NOT NULL;DEFAULT:0"`
}

// MigrationTasks 注册数据迁移任务
func MigrationTasks() []*migratorx.Task {
	return []*migratorx.Task{
		{Comment: "用户表", Model: &User{}},
	}
}
