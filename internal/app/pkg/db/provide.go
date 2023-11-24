package db

import (
	"context"
	"database/sql"

	"go-scaffold/internal/config"
)

// Provide database connection
func Provide(ctx context.Context, conf config.DatabaseConn) (db *sql.DB, cleanup func(), err error) {
	db, err = New(ctx, conf)
	if err != nil {
		return nil, nil, err
	}

	cleanup = func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}

	return
}
