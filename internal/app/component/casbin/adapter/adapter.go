package adapter

import (
	fileadapter "go-scaffold/internal/app/component/casbin/adapter/file"
	gormadapter "go-scaffold/internal/app/component/casbin/adapter/gorm"

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

type Config struct {
	File *fileadapter.Config
	Gorm *gormadapter.Config
}

func New(config *Config, db *gorm.DB) (adp Adapter, err error) {
	if config == nil {
		return nil, nil
	}

	if config.Gorm != nil {
		adp, err = gormadapter.New(config.Gorm, db)
		if err != nil {
			return nil, err
		}
	}

	if config.File != nil {
		adp = fileadapter.New(config.File)
	}

	return
}
