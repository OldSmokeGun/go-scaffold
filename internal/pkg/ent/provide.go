package ent

import (
	"database/sql"
	"log/slog"

	"go-scaffold/internal/config"
	"go-scaffold/internal/pkg/ent/ent"
)

// Provide db client
func Provide(env config.Env, conf config.DatabaseConn, logger *slog.Logger, db *sql.DB) (*ent.Client, func(), error) {
	client, err := New(env, conf, logger, db)
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
