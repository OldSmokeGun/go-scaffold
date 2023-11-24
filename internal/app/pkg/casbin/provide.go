package casbin

import (
	"database/sql"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go-scaffold/internal/config"
)

// Provide casbin
func Provide(
	env config.Env,
	conf config.Casbin,
	dbConf config.DatabaseConn,
	logger *slog.Logger,
	gdb *gorm.DB,
	sdb *sql.DB,
) (*casbin.Enforcer, error) {
	return New(env, conf, dbConf, logger, gdb, sdb)
}
