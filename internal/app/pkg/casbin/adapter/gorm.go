package adapter

import (
	"go-scaffold/internal/config"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewGormAdapter build casin gorm adapter
func NewGormAdapter(conf config.CasbinGormAdapter, db *gorm.DB) (adp *gormadapter.Adapter, err error) {
	gormadapter.TurnOffAutoMigrate(db)

	if conf.TableName == "" {
		adp, err = gormadapter.NewAdapterByDB(db)
		if err != nil {
			return nil, err
		}
	} else {
		adp, err = gormadapter.NewAdapterByDBUseTableName(db, "", conf.TableName)
		if err != nil {
			return nil, err
		}
	}

	return
}
