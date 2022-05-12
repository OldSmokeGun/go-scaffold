package gorm

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const defaultTableName = "casbin_rule"

type Config struct {
	TableName string
}

// New casin gorm adapter
func New(config *Config, db *gorm.DB) (adp *gormadapter.Adapter, err error) {
	tableName := defaultTableName

	if config.TableName != "" {
		tableName = config.TableName
	}

	adp, err = gormadapter.NewAdapterByDBUseTableName(db, "", tableName)
	if err != nil {
		return nil, err
	}

	return
}
