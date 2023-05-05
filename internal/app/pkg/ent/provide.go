package ent

import (
	"context"

	"go-scaffold/internal/app/pkg/ent/ent"
	"go-scaffold/internal/config"

	"golang.org/x/exp/slog"
)

// Provide db client
func Provide(ctx context.Context, env config.Env, conf config.DBConn, logger *slog.Logger) (*ent.Client, func(), error) {
	client, err := New(ctx, env, conf, logger)
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
