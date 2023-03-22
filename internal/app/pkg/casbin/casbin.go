package casbin

import (
	"go-scaffold/internal/app/pkg/casbin/adapter"
	"go-scaffold/internal/app/pkg/casbin/model"
	"go-scaffold/internal/config"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// New build casbin
func New(conf config.Casbin, db *gorm.DB) (*casbin.Enforcer, error) {
	mod, err := model.New(conf.Model)
	if err != nil {
		return nil, err
	}

	adp, err := adapter.New(conf.Adapter, db)
	if err != nil {
		return nil, err
	}

	ef, err := casbin.NewEnforcer(mod, adp)
	if err != nil {
		return nil, err
	}

	return ef, nil
}
