package casbin

import (
	"database/sql"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go-scaffold/internal/app/pkg/casbin/adapter"
	"go-scaffold/internal/app/pkg/casbin/model"
	"go-scaffold/internal/config"
)

// New build casbin
func New(
	env config.Env,
	conf config.Casbin,
	dbConf config.DatabaseConn,
	logger *slog.Logger,
	gdb *gorm.DB,
	sdb *sql.DB,
) (*casbin.Enforcer, error) {
	mod, err := model.New(conf.Model)
	if err != nil {
		return nil, err
	}

	adp, err := adapter.New(env, conf.Adapter, dbConf, logger, gdb, sdb)
	if err != nil {
		return nil, err
	}

	ef, err := casbin.NewEnforcer(mod, adp)
	if err != nil {
		return nil, err
	}

	return ef, nil
}
