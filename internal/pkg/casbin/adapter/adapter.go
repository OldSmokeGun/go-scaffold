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

// New creates casbin adapter.
// A configured gorm adapter is used even when a file path is also set.
// The file adapter does not implement AddPolicy, so rules would stay in memory
// and never reach casbin_rules.
func New(
	conf config.CasbinAdapter,
	db *gorm.DB,
) (Adapter, error) {
	if conf.Gorm != nil {
		return NewGormAdapter(db)
	}

	if conf.File != "" {
		return NewFileAdapter(conf.File), nil
	}

	return nil, nil
}
