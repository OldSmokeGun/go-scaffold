package schema

import (
	"go-scaffold/internal/app/domain"
	igorm "go-scaffold/internal/pkg/gorm"
)

type SystemRole struct {
	igorm.BaseModel `gorm:"embedded"`
	Name            string `gorm:"column:name;uniqueIndex;size:32;not null;default:''"`
}

func (SystemRole) TableName() string {
	return "system_roles"
}

func (m *SystemRole) ToEntity() *domain.SystemRole {
	return &domain.SystemRole{
		ID:   m.ID,
		Name: m.Name,
	}
}
