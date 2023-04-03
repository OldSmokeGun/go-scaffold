package repository

import (
	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(domain.UserRepository), new(*UserRepository)), NewUserRepository),
)

// BaseModel 基础模型
// 自动更新时间戳、软删除
type BaseModel struct {
	ID        int64                 `gorm:"primaryKey"`
	CreatedAt int64                 `gorm:"NOT NULL"`
	UpdatedAt int64                 `gorm:"NOT NULL;DEFAULT:0"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;NOT NULL;DEFAULT:0"`
}

// convertError 转换 gorm 包的错误为内部错误
//
// 屏蔽 repository 层的内部实现
func convertError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(berr.ErrResourceNotFound)
	}
	return errors.WithStack(err)
}
