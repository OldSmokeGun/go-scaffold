package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"go-scaffold/internal/config"
)

// ErrUnsupportedDriver unsupported database driver
var ErrUnsupportedDriver = errors.New("unsupported database driver")

// New build database connection
func New(ctx context.Context, conf config.DBConn) (*sql.DB, error) {
	if !conf.Driver.IsSupported() {
		return nil, ErrUnsupportedDriver
	}

	dsn := buildDSN(conf)

	db, err := sql.Open(conf.Driver.String(), dsn)
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
		db.SetConnMaxLifetime(conf.ConnMaxIdleTime)
	}

	if conf.ConnMaxLifeTime > 0 {
		db.SetConnMaxLifetime(conf.ConnMaxLifeTime)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func buildDSN(c config.DBConn) string {
	dsn := ""

	switch c.Driver {
	case config.Postgres:
		var host, port string
		s := strings.SplitN(c.Addr, ":", 2)
		if len(s) == 2 {
			host = s[0]
			port = s[1]
		}

		dsn = "host=" + host + " port=" + port + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " " + c.Options
	case config.MySQL:
		fallthrough
	default:
		dsn = c.Username + ":" + c.Password + "@tcp(" + c.Addr + ")/" + c.Database + "?" + c.Options
	}

	return dsn
}
