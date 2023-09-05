package ent

import (
	"database/sql"
	"log/slog"

	"go-scaffold/internal/app/pkg/ent/ent"
	"go-scaffold/internal/config"
)

// Provide db client
func Provide(env config.Env, conf config.DBConn, logger *slog.Logger, db *sql.DB) (*ent.Client, func(), error) {
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
