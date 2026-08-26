package schema

import (
	"go-scaffold/internal/app/domain"
	igorm "go-scaffold/internal/pkg/gorm"
)

type Product struct {
	igorm.BaseModel `gorm:"embedded"`
	Name            string `gorm:"column:name;index;size:128;not null;default:''"`
	Desc            string `gorm:"column:desc;size:255;not null;default:''"`
	Price           int    `gorm:"column:price;not null;default:0"`
}

func (Product) TableName() string {
	return "products"
}

func (m *Product) ToEntity() *domain.Product {
	return &domain.Product{
		ID:    m.ID,
		Name:  m.Name,
		Desc:  m.Desc,
		Price: m.Price,
	}
}
