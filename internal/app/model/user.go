package model

// User 用户表
type User struct {
	BaseModel
	Name string `gorm:"column:name;type:varchar(64);not null;default:'';comment:用户名"`
}

func (receiver User) TableName() string {
	return "user"
}
