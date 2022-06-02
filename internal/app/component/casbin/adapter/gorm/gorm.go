package gorm

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Config struct {
	TableName string

	migration func(db *gorm.DB) error
}

func (c *Config) SetMigration(fn func(db *gorm.DB) error) {
	c.migration = fn
}

// New casin gorm adapter
func New(config *Config, db *gorm.DB) (adp *gormadapter.Adapter, err error) {
	if config.TableName == "" {
		adp, err = gormadapter.NewAdapterByDB(db)
		if err != nil {
			return nil, err
		}
	} else {
		db = gormadapter.TurnOffAutoMigrate(db)

		if config.migration != nil {
			if err = config.migration(db); err != nil {
				return nil, err
			}
		}

		adp, err = gormadapter.NewAdapterByDBUseTableName(db, "", config.TableName)
		if err != nil {
			return nil, err
		}
	}

	return
}
