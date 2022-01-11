package model

import (
	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID        uint                  `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt int                   `json:"created_at,omitempty"`
	UpdatedAt int                   `json:"updated_at,omitempty"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
