package adapter

import (
	"github.com/casbin/casbin/v2/persist"
	"gorm.io/gorm"

	"go-scaffold/internal/config"
)

// Adapter the interface that casbin adapter must implement
type Adapter interface {
	persist.Adapter
	persist.BatchAdapter
	persist.UpdatableAdapter
	persist.FilteredAdapter
}

// New creates casbin adapter
func New(
	conf config.CasbinAdapter,
	db *gorm.DB,
) (adp Adapter, err error) {
	if conf.Gorm != nil {
		adp, err = NewGormAdapter(db)
		if err != nil {
			return nil, err
		}
	}

	if conf.File != "" {
		adp = NewFileAdapter(conf.File)
	}

	return
}
