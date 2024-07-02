package ent

import (
	"context"
	"log/slog"

	"go-scaffold/internal/config"
	"go-scaffold/internal/pkg/ent/ent"
)

type DefaultClient = ent.Client

// ProvideDefault db client
func ProvideDefault(ctx context.Context, env config.Env, conf config.DefaultDatabase, logger *slog.Logger) (*DefaultClient, func(), error) {
	client, err := New(ctx, env, conf.DatabaseConn, logger)
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
