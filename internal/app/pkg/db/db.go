package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx"
	_ "github.com/mattn/go-sqlite3"

	"go-scaffold/internal/config"
)

// ErrUnsupportedDriver unsupported database driver
var ErrUnsupportedDriver = errors.New("unsupported database driver")

// New build database connection
func New(ctx context.Context, conf config.DBConn) (*sql.DB, error) {
	if !conf.Driver.IsSupported() {
		return nil, ErrUnsupportedDriver
	}

	db, err := sql.Open(conf.Driver.String(), conf.DSN)
	if err != nil {
		return nil, err
	}

	if conf.MaxIdleConn > 0 {
		db.SetMaxIdleConns(conf.MaxIdleConn)
	}

	if conf.MaxOpenConn > 0 {
		db.SetMaxOpenConns(conf.MaxOpenConn)
	}

	if conf.ConnMaxIdleTime > 0 {
		db.SetConnMaxLifetime(conf.ConnMaxIdleTime * time.Second)
	}

	if conf.ConnMaxLifeTime > 0 {
		db.SetConnMaxLifetime(conf.ConnMaxLifeTime * time.Second)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
