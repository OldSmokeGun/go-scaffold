package model

import (
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	BaseModel
	Name  string `gorm:"column:name;type:varchar(64);not null;default:'';comment:名称"`
	Age   int8   `gorm:"column:age;type:tinyint(3);not null;default:0;comment:年龄"`
	Phone string `gorm:"column:phone;type:varchar(11);not null;default:'';comment:手机号码"`
}

func (u User) TableName() string {
	return "user"
}

func (u User) Migrate(db *gorm.DB) error {
	if err := db.Set(
		"gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表'",
	).AutoMigrate(u); err != nil {
		return err
	}

	return nil
}
