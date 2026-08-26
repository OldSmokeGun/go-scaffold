package casbin

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go-scaffold/internal/config"
	"go-scaffold/internal/pkg/casbin/adapter"
	"go-scaffold/internal/pkg/casbin/model"
)

// New build casbin
func New(conf config.Casbin, gdb *gorm.DB) (*casbin.Enforcer, error) {
	mod, err := model.New(conf.Model)
	if err != nil {
		return nil, err
	}

	adp, err := adapter.New(conf.Adapter, gdb)
	if err != nil {
		return nil, err
	}

	return casbin.NewEnforcer(mod, adp)
}
