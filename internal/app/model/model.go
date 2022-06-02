package model

import (
	"errors"
	"go-scaffold/internal/app/model/migrator"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	ErrDataStoreFailed  = errors.New("数据保存失败")
	ErrDataQueryFailed  = errors.New("数据查询失败")
	ErrDataDeleteFailed = errors.New("数据删除失败")
	ErrDataNotFound     = errors.New("数据不存在")
)

// BaseModel 基础模型
// 自动更新时间戳、软删除
type BaseModel struct {
	Id        uint64                `gorm:"primaryKey"`
	CreatedAt int64                 `gorm:"NOT NULL"`
	UpdatedAt int64                 `gorm:"NOT NULL;DEFAULT:0"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;NOT NULL;DEFAULT:0"`
}

// Migrate 注册数据迁移模型
func Migrate(db *gorm.DB) error {
	models := []migrator.Migrator{
		&User{},
	}

	for _, model := range models {
		if err := model.Migrate(db); err != nil {
			return err
		}
	}

	return nil
}
