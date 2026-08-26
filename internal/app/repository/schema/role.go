package schema

import (
	"go-scaffold/internal/app/domain"
	igorm "go-scaffold/internal/pkg/gorm"
)

type Role struct {
	igorm.BaseModel `gorm:"embedded"`
	Name            string `gorm:"column:name;uniqueIndex;size:32;not null;default:''"`
}

func (Role) TableName() string {
	return "roles"
}

func (m *Role) ToEntity() *domain.Role {
	return &domain.Role{
		ID:   m.ID,
		Name: m.Name,
	}
}
