package casbin

import (
	"context"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go-scaffold/internal/config"
	"go-scaffold/internal/pkg/db"
)

// Provide casbin
func Provide(
	ctx context.Context,
	env config.Env,
	conf config.Casbin,
	dbConf config.DefaultDatabase,
	logger *slog.Logger,
	gdb *gorm.DB,
) (*casbin.Enforcer, error) {
	sdb, err := db.New(ctx, dbConf.DatabaseConn)
	if err != nil {
		return nil, err
	}

	return New(env, conf, dbConf.DatabaseConn, logger, gdb, sdb)
}
