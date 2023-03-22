package gorm

import (
	"go-scaffold/internal/app/pkg/db"
	"go-scaffold/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDialect build dialect
func NewDialect(driver config.DBDriver, conn gorm.ConnPool) (gorm.Dialector, error) {
	var dialect gorm.Dialector
	switch driver {
	case config.MySQL:
		dialect = newMySQLDialect(conn)
	case config.Postgres:
		dialect = newPostgresDialect(conn)
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
