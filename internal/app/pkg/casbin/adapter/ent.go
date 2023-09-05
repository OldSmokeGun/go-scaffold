package adapter

import (
	"database/sql"
	"log/slog"

	entsql "entgo.io/ent/dialect/sql"
	entadapter "github.com/casbin/ent-adapter"
	"github.com/casbin/ent-adapter/ent"

	"go-scaffold/internal/config"
	elog "go-scaffold/pkg/log/ent"
)

// NewEntAdapter build casin ent adapter
func NewEntAdapter(env config.Env, dbConf config.DBConn, logger *slog.Logger, sdb *sql.DB) (*entadapter.Adapter, error) {
	driver := entsql.OpenDB(dbConf.Driver.String(), sdb)

	options := []ent.Option{
		ent.Driver(driver),
		ent.Log(elog.NewLogger(logger).Log),
	}
	if env.IsDebug() {
		options = append(options, ent.Debug())
	}

	client := ent.NewClient(options...)

	adapter, err := entadapter.NewAdapterWithClient(client)
	if err != nil {
		return nil, err
	}

	return adapter, nil
}
