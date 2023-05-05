package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/execquery --feature sql/modifier --feature intercept --target ./ent ../../repository/schema

import (
	"context"

	"go-scaffold/internal/app/pkg/db"
	"go-scaffold/internal/app/pkg/ent/ent"
	_ "go-scaffold/internal/app/pkg/ent/ent/runtime"
	"go-scaffold/internal/config"
	elog "go-scaffold/pkg/log/ent"

	entsql "entgo.io/ent/dialect/sql"
	"golang.org/x/exp/slog"
)

// New build db client
func New(ctx context.Context, env config.Env, dbConf config.DBConn, logger *slog.Logger) (*ent.Client, error) {

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
