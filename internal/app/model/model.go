package model

import (
	"go-scaffold/internal/app/pkg/migrator"
	"gorm.io/plugin/soft_delete"
)

// BaseModel 基础模型
// 自动更新时间戳、软删除
type BaseModel struct {
	ID        uint                  `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt int                   `json:"created_at,omitempty"`
	UpdatedAt int                   `json:"updated_at,omitempty"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// MigrationTasks 注册数据迁移任务
func MigrationTasks() []*migrator.Task {
	return []*migrator.Task{
		{Comment: "用户表", Model: &User{}},
	}
}
