package models

import "gorm.io/gorm"

type BaseModel struct {
	ID        uint `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt int  `json:"created_at,omitempty"`
	UpdatedAt int  `json:"updated_at,omitempty"`
	DeletedAt int  `json:"deleted_at,omitempty"`
}

type Pagination struct {
	Page     int
	PageSize int
}

func PaginationScope(p Pagination) func(db *gorm.DB) *gorm.DB {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PageSize <= 0 {
		p.PageSize = 10
	}

	offset := (p.Page - 1) * p.PageSize

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(p.PageSize)
	}
}

func StatusEnableScope(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", 1)
}
