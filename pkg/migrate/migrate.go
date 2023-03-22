package migrate

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Driver database driver
type Driver string

const (
	MySQL    Driver = "mysql"
	Postgres Driver = "postgres"
)

// New build migrator
func New(source, databaseURL string) (*migrate.Migrate, error) {
	return migrate.New(source, databaseURL)
}

// NewWithDB build migrator with db
func NewWithDB(source string, driver Driver, db *sql.DB) (*migrate.Migrate, error) {
	var (
		instance database.Driver
		err      error
	)

	switch driver {
	case MySQL:
		instance, err = mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return nil, err
		}
	case Postgres:
		instance, err = postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
	}

	return migrate.NewWithDatabaseInstance(source, string(driver), instance)
}
