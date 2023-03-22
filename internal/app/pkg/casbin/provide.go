package casbin

import (
	"go-scaffold/internal/config"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// Provide casbin
func Provide(conf config.Casbin, db *gorm.DB) (*casbin.Enforcer, error) {
	return New(conf, db)
}
