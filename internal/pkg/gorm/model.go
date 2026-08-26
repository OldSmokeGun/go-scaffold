package gorm

import "gorm.io/plugin/soft_delete"

type BaseModel struct {
	ID        int64                 `gorm:"column:id;primaryKey"`
	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime;not null;default:0"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;index;not null;default:0" gen:"softDelete"`
}
