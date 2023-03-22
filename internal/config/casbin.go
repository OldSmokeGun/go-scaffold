package config

// Casbin casbin config
type Casbin struct {
	Model   CasbinModel   `json:"model"`
	Adapter CasbinAdapter `json:"adapter"`
}

func (Casbin) GetName() string {
	return "casbin"
}

// CasbinModel casbin model
type CasbinModel struct {
	Path string `json:"path"`
}

// CasbinFileAdapter casbin file adapter
type (
	// CasbinAdapter casbin adapter
	CasbinAdapter struct {
		File *CasbinFileAdapter `json:"file"`
		Gorm *CasbinGormAdapter `json:"gorm"`
	}

	CasbinFileAdapter struct {
		Path string `json:"path"`
	}

	// CasbinGormAdapter casbin gorm adapter
	CasbinGormAdapter struct {
		TableName string `json:"tableName"`
	}
)
