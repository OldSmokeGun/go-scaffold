package model

// User 用户表
type User struct {
	BaseModel
	Name  string `gorm:"column:name;type:varchar(64);not null;default:'';comment:名称"`
	Age   int8   `gorm:"column:age;type:tinyint(3);not null;default:0;comment:年龄"`
	Phone string `gorm:"column:phone;type:varchar(11);not null;default:'';comment:手机号码"`
}

func (receiver User) TableName() string {
	return "user"
}
