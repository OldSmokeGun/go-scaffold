package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/execquery --target ./ent ../../repository/schema

import (
	"context"

	"go-scaffold/internal/app/pkg/db"
	"go-scaffold/internal/app/pkg/ent/ent"
	"go-scaffold/internal/config"

	entsql "entgo.io/ent/dialect/sql"
	"golang.org/x/exp/slog"
)

// New build db client
func New(ctx context.Context, conf config.DBConn, logger *slog.Logger) (*ent.Client, error) {
	sdb, err := db.New(ctx, conf)
	if err != nil {
		return nil, err
	}

	driver := entsql.OpenDB(conf.Driver.String(), sdb)
	client := ent.NewClient(
		ent.Driver(driver),
		ent.Log(func(i ...any) {
			if logger != nil {
				logger.Debug("ent debug", i...)
			}
		}),
	)

	return client, nil
}
