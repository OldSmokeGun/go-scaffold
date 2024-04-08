package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-scaffold/internal/config"
	"go-scaffold/internal/pkg/db"
)

// NewDialect build dialect
func NewDialect(driver config.DatabaseDriver, conn gorm.ConnPool) (gorm.Dialector, error) {
	var dialect gorm.Dialector
	switch driver {
	case config.MySQL:
		dialect = newMySQLDialect(conn)
	case config.Postgres:
		dialect = newPostgresDialect(conn)
	case config.SQLite:
		dialect = newSQLiteDialect(conn)
	default:
		return nil, db.ErrUnsupportedDriver
	}
	return dialect, nil
}

// newMySQLDialect build mysql dialect
func newMySQLDialect(conn gorm.ConnPool) gorm.Dialector {
	return mysql.New(mysql.Config{Conn: conn})
}

// newPostgresDialect build postgres dialect
func newPostgresDialect(conn gorm.ConnPool) gorm.Dialector {
	return postgres.New(postgres.Config{Conn: conn})
}

// newSQLiteDialect build sqlite dialect
func newSQLiteDialect(conn gorm.ConnPool) gorm.Dialector {
	return sqlite.Dialector{Conn: conn}
}
