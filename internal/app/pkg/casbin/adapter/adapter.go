package adapter

import (
	"go-scaffold/internal/config"

	"github.com/casbin/casbin/v2/persist"
	"gorm.io/gorm"
)

// Adapter the interface that casbin adapter must implement
type Adapter interface {
	persist.Adapter
	persist.BatchAdapter
	persist.UpdatableAdapter
	persist.FilteredAdapter
}

// New creates casin adapter
func New(conf config.CasbinAdapter, db *gorm.DB) (adp Adapter, err error) {
	if conf.Gorm != nil {
		adp, err = NewGormAdapter(*conf.Gorm, db)
		if err != nil {
			return nil, err
		}
	}

	if conf.File != nil {
		adp = NewFileAdapter(*conf.File)
	}

	return
}
