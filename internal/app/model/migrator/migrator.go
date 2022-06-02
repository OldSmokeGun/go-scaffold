package migrator

import "gorm.io/gorm"

// Migrator 模型数据迁移接口
type Migrator interface {
	Migrate(db *gorm.DB) error
}
