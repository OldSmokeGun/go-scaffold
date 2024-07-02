package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/execquery --feature sql/modifier --feature intercept --target ./ent ../../app/repository/schema

import (
	"context"
	"log/slog"

	entsql "entgo.io/ent/dialect/sql"

	"go-scaffold/internal/config"
	"go-scaffold/internal/pkg/db"
	"go-scaffold/internal/pkg/ent/ent"
	_ "go-scaffold/internal/pkg/ent/ent/runtime"
	elog "go-scaffold/pkg/log/ent"
)

// New build db client
func New(ctx context.Context, env config.Env, dbConf config.DatabaseConn, logger *slog.Logger) (*ent.Client, error) {
	sdb, err := db.New(ctx, dbConf)
	if err != nil {
		return nil, err
	}

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
