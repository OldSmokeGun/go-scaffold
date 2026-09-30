package schema

import (
	"go-scaffold/internal/app/domain"
	igorm "go-scaffold/internal/pkg/gorm"
)

type SystemUser struct {
	igorm.BaseModel `gorm:"embedded"`
	Username        string `gorm:"column:username;uniqueIndex;size:32;not null;default:''"`
	Password        string `gorm:"column:password;size:64;not null;default:''"`
	Nickname        string `gorm:"column:nickname;size:64;not null;default:''"`
	Phone           string `gorm:"column:phone;index;size:11;not null;default:''"`
	Salt            string `gorm:"column:salt;size:64;not null;default:''"`
}

func (SystemUser) TableName() string {
	return "system_users"
}

func (m *SystemUser) ToEntity() *domain.SystemUser {
	return &domain.SystemUser{
		ID:       m.ID,
		Username: m.Username,
		Password: domain.Password(m.Password),
		Nickname: m.Nickname,
		Phone:    m.Phone,
		Salt:     m.Salt,
	}
}
