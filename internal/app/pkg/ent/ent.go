package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/execquery --feature sql/modifier --feature intercept --target ./ent ../../repository/schema

import (
	"database/sql"
	"log/slog"

	entsql "entgo.io/ent/dialect/sql"

	"go-scaffold/internal/app/pkg/ent/ent"
	_ "go-scaffold/internal/app/pkg/ent/ent/runtime"
	"go-scaffold/internal/config"
	elog "go-scaffold/pkg/log/ent"
)

// New build db client
func New(env config.Env, dbConf config.DatabaseConn, logger *slog.Logger, sdb *sql.DB) (*ent.Client, error) {
	driver := entsql.OpenDB(dbConf.Driver.String(), sdb)

	options := []ent.Option{
		ent.Driver(driver),
		ent.Log(elog.NewLogger(logger).Log),
	}
	if env.IsDebug() {
		options = append(options, ent.Debug())
	}

	client := ent.NewClient(options...)

	return client, nil
}
