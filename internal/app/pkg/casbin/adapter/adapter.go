package adapter

import (
	"database/sql"
	"log/slog"

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

// New creates casin adapter
func New(
	env config.Env,
	conf config.CasbinAdapter,
	dbConf config.DatabaseConn,
	logger *slog.Logger,
	db *gorm.DB,
	sdb *sql.DB,
) (adp Adapter, err error) {
	if conf.Gorm != nil {
		adp, err = NewGormAdapter(db)
		if err != nil {
			return nil, err
		}
	}

	if conf.Ent != nil {
		adp, err = NewEntAdapter(env, dbConf, logger, sdb)
		if err != nil {
			return nil, err
		}
	}

	if conf.File != "" {
		adp = NewFileAdapter(conf.File)
	}

	return
}
