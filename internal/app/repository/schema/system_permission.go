package schema

import (
	"go-scaffold/internal/app/domain"
	igorm "go-scaffold/internal/pkg/gorm"
)

type SystemPermission struct {
	igorm.BaseModel `gorm:"embedded"`
	Key             string `gorm:"column:key;uniqueIndex;size:128;not null;default:''"`
	Name            string `gorm:"column:name;size:128;not null;default:''"`
	Desc            string `gorm:"column:desc;size:255;not null;default:''"`
	ParentID        int64  `gorm:"column:parent_id;not null;default:0"`
}

func (SystemPermission) TableName() string {
	return "system_permissions"
}

func (m *SystemPermission) ToEntity() *domain.SystemPermission {
	return &domain.SystemPermission{
		ID:       m.ID,
		Key:      m.Key,
		Name:     m.Name,
		Desc:     m.Desc,
		ParentID: m.ParentID,
	}
}
