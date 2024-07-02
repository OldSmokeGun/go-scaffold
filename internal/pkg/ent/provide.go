package ent

import (
	"context"
	"log/slog"

	"go-scaffold/internal/config"
	"go-scaffold/internal/pkg/db"
	"go-scaffold/internal/pkg/ent/ent"
)

type DefaultClient = ent.Client

// ProvideDefault db client
func ProvideDefault(ctx context.Context, env config.Env, conf config.DefaultDatabase, logger *slog.Logger) (*DefaultClient, func(), error) {
	sdb, err := db.New(ctx, conf.DatabaseConn)
	if err != nil {
		return nil, nil, err
	}

	client, err := New(env, conf.DatabaseConn, logger, sdb)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		err = client.Close()
		if err != nil {
			panic(err)
		}
	}

	return client, cleanup, nil
}
