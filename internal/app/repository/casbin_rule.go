package repository

import "gorm.io/gorm"

type casbinRuleModel struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100"`
	V0    string `gorm:"size:100"`
	V1    string `gorm:"size:100"`
	V2    string `gorm:"size:100"`
	V3    string `gorm:"size:100"`
	V4    string `gorm:"size:100"`
	V5    string `gorm:"size:100"`
	V6    string `gorm:"size:25"`
	V7    string `gorm:"size:25"`
}

// TableName 表名
func (r casbinRuleModel) TableName() string {
	return "casbin_rules"
}

// Migrate 迁移
func (r casbinRuleModel) Migrate(db *gorm.DB) error {
	if err := db.Set(
		"gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表'",
	).AutoMigrate(r); err != nil {
		return err
	}

	return nil
}
