package casbin

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go-scaffold/internal/config"
)

// Provide casbin
func Provide(conf config.Casbin, gdb *gorm.DB) (*casbin.Enforcer, error) {
	return New(conf, gdb)
}
