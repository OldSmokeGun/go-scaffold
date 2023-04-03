package repository

type casbinRuleModel struct {
	ID    uint   `gorm:"column:id;primaryKey"`
	Ptype string `gorm:"column:ptype"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
	V3    string `gorm:"column:v3"`
	V4    string `gorm:"column:v4"`
	V5    string `gorm:"column:v5"`
}

// TableName 表名
func (r casbinRuleModel) TableName() string {
	return "casbin_rules"
}
